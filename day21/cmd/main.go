package main

import (
	"fmt"
	"io"
	"regexp"
	"strings"

	common "github.com/torbensky/adventofcode-common"
)

type stringSet map[string]struct{}

func (s stringSet) only() string {
	if len(s) != 1 {
		panic("set does not have only 1")
	}
	for k := range s {
		return k
	}

	panic("set does not have only 1")
}

// returns a set that is the intersection of sets a and b
func intersect(a, b stringSet) stringSet {
	result := make(stringSet)

	for k := range a {
		if _, ok := b[k]; ok {
			result[k] = struct{}{}
		}
	}

	for k := range b {
		if _, ok := a[k]; ok {
			result[k] = struct{}{}
		}
	}

	return result
}

func diff(a, b stringSet) stringSet {
	result := make(stringSet)
	for k := range a {
		if _, ok := b[k]; !ok {
			result[k] = struct{}{}
		}
	}
	return result
}

func (s stringSet) union(o stringSet) {
	for k := range o {
		s[k] = struct{}{}
	}
}

func fromSlice(slice []string) stringSet {
	result := make(stringSet)
	for _, s := range slice {
		result[s] = struct{}{}
	}
	return result
}

func main() {
	fmt.Printf("Part 1: %d\n", part1(common.OpenInputFile()))
	fmt.Printf("Part 2: %s\n", part2(common.OpenInputFile()))
}

type foodList struct {
	ingredients []string
	allergens   []string
}

var lineRegex = regexp.MustCompile("(.+?)\\(contains(.+)\\)")

func parseFoodList(reader io.Reader) (ingredientCounts map[string]int, allIngredients stringSet, allergenToPossibleIng map[string]stringSet) {
	var foodLists []foodList
	ingredientCounts = make(map[string]int)
	allIngredients = make(stringSet)
	common.ScanLines(reader, func(line string) {
		matches := lineRegex.FindStringSubmatch(line)
		if len(matches) != 3 {
			panic("Not enough matches")
		}

		// Add to count of ingredients
		ingredients := strings.Split(strings.TrimSpace(matches[1]), " ")
		for _, i := range ingredients {
			ingredientCounts[i]++
			allIngredients[i] = struct{}{}
		}
		// fmt.Println(ingredients)

		allergens := strings.Split(strings.TrimSpace(matches[2]), ", ")
		// fmt.Println(allergens)
		foodLists = append(foodLists, foodList{
			ingredients: ingredients,
			allergens:   allergens,
		})
	})

	// Map all the allergens to their possible ingredient sources
	allergenToPossibleIng = make(map[string]stringSet)
	for _, fl := range foodLists {
		for _, a := range fl.allergens {
			cur, ok := allergenToPossibleIng[a]
			if ok {
				// We can narrow down what we know so far with this new info
				allergenToPossibleIng[a] = intersect(cur, fromSlice(fl.ingredients))
			} else {
				// First occurrence of the allergen
				allergenToPossibleIng[a] = fromSlice(fl.ingredients)
			}
		}
	}

	return ingredientCounts, allIngredients, allergenToPossibleIng
}

func part1(reader io.Reader) int {
	ingredientCounts, allIngredients, allergenToPossibleIng := parseFoodList(reader)
	mightHaveAllergen := make(stringSet)
	for _, p := range allergenToPossibleIng {
		mightHaveAllergen.union(p)
	}

	// fmt.Println("maybe:", mightHaveAllergen)
	// fmt.Println("all:", allIngredients)

	inert := diff(allIngredients, mightHaveAllergen)
	// fmt.Println("inert:", inert)

	total := 0
	for i := range inert {
		total += ingredientCounts[i]
	}

	return total
}

func part2(reader io.Reader) string {
	_, _, allergenToPossibleIng := parseFoodList(reader)

	clearIdentified := func(ingred string) {
		for _, il := range allergenToPossibleIng {
			delete(il, ingred)
		}
	}

	identified := make(map[string]string)
	for {
		for a, il := range allergenToPossibleIng {
			if len(il) == 1 {
				i := il.only()
				identified[a] = i
				delete(allergenToPossibleIng, a)
				clearIdentified(i)
			}
		}

		if len(allergenToPossibleIng) == 0 {
			break
		}
	}

	return fmt.Sprint(identified)
}
