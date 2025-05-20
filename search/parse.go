/*=================================== LICENSE =======================================

                                   Apache License
                             Version 2.0, January 2004
                          http://www.apache.org/licenses/

============================== BRIEF FILE DESCRIPTION ===============================

This file contains the implementation of the ParseQuery function, which is responsible
for translating a sequence of fragments into a well-formed query (a tree-like structure).

The parsing itself is done using a stack-based approach, where operators and operands
are pushed onto their respective stacks. The precedence of operators is taken into
account to ensure that the resulting query structure is correct.
==================================================================================*/

package search

import (
	"fmt"
	"quinto/core"
)

type parsingState struct {
	queryStack      *[]core.Query
	opStack         *[]any
	precedenceStack *[]int
}

type openParenthesis struct {
	openingPosition int
}

func stackPush[T any](stack *[]T, q T) {
	*stack = append(*stack, q)
}

func stackPop[T any](stack *[]T) T {
	if len(*stack) == 0 {
		var zero T
		return zero
	}
	q := (*stack)[len(*stack)-1]
	*stack = (*stack)[:len(*stack)-1]
	return q
}

func evaluateOne(queryStack *[]core.Query, opStack *[]any) {
	op := stackPop(opStack)
	switch v := op.(type) {
	case ComplexQuery:
		v.rx = stackPop(queryStack)
		v.lx = stackPop(queryStack)
		var castedToQuery core.Query = &v
		stackPush(queryStack, castedToQuery)
	default:
		return
	}
}

func evaluateAll(state parsingState, currentPrecedence int) {
	counter := len(*state.precedenceStack) - 1
	for counter >= 0 && (*state.precedenceStack)[counter] >= currentPrecedence {
		evaluateOne(state.queryStack, state.opStack)
		stackPop(state.precedenceStack)
		counter--
	}
}

func ParseQuery(queryFragments []queryFragment) (core.Query, error) {

	var queryStack []core.Query
	var opStack []any
	var precedenceStack []int
	var parsingState = parsingState{
		queryStack:      &queryStack,
		opStack:         &opStack,
		precedenceStack: &precedenceStack,
	}

	for index, fragment := range queryFragments {
		switch fragment.txt {
		case "(":
			var castedAsAny any = openParenthesis{openingPosition: index}
			stackPush(&opStack, castedAsAny)
			stackPush(&precedenceStack, 0)
		case "OR":
			evaluateAll(parsingState, 1)
			var castedAsAny any = ComplexQuery{ord: fragment.ord, policy: OrQueryPolicy}
			stackPush(&opStack, castedAsAny)
			stackPush(&precedenceStack, 1)
		case "XOR":
			evaluateAll(parsingState, 2)
			var castedAsAny any = ComplexQuery{ord: fragment.ord, policy: XorQueryPolicy}
			stackPush(&opStack, castedAsAny)
			stackPush(&precedenceStack, 2)
		case "AND":
			evaluateAll(parsingState, 3)
			var castedAsAny any = ComplexQuery{ord: fragment.ord, policy: AndQueryPolicy}
			stackPush(&opStack, castedAsAny)
			stackPush(&precedenceStack, 3)
		case "NEAR":
			evaluateAll(parsingState, 4)
			var castedAsAny any = ComplexQuery{ord: fragment.ord, policy: NearQueryPolicy(fragment.opt)}
			stackPush(&opStack, castedAsAny)
			stackPush(&precedenceStack, 4)
		case ")":
			for len(opStack) > 0 {
				for precedence := stackPop(&precedenceStack); precedence != 0; {
					evaluateOne(&queryStack, &opStack)
					precedence = stackPop(&precedenceStack)
				}
				stackPop(&opStack)
			}
		default:
			queryStack = append(queryStack, &ExactQuery{term: fragment.txt})
		}
	}

	for len(opStack) > 0 {
		evaluateOne(&queryStack, &opStack)
	}

	if len(queryStack) != 1 {
		return nil, fmt.Errorf("invalid query: %v", queryStack)
	}

	return queryStack[0], nil
}
