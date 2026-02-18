package converter

import (
	"testing"
)

// Benchmark for ConvertDiceRolls - tests impact of regex compilation
func BenchmarkConvertDiceRolls(b *testing.B) {
	text := "Attack: to hit: 1d20+5, damage: 2d6+3, and healing: 1d8+2"
	actionName := "Longsword"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ConvertDiceRolls(text, actionName)
	}
}

// Benchmark for extractModifier - tests impact of regex compilation
func BenchmarkExtractModifier(b *testing.B) {
	notation := "1d20+5"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = extractModifier(notation)
	}
}

// Benchmark for isD20Roll - tests impact of regex compilation
func BenchmarkIsD20Roll(b *testing.B) {
	notation := "1d20+5"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = isD20Roll(notation)
	}
}

// Benchmark for calculateAverage - tests impact of regex compilation
func BenchmarkCalculateAverage(b *testing.B) {
	notation := "2d6+3"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = calculateAverage(notation)
	}
}

// Benchmark for ParseDiceNotation - tests impact of regex compilation
func BenchmarkParseDiceNotation(b *testing.B) {
	notation := "2d6+3"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ParseDiceNotation(notation)
	}
}

// Benchmark with multiple different dice notations
func BenchmarkConvertDiceRollsVaried(b *testing.B) {
	texts := []string{
		"Attack: to hit: 1d20+5, damage: 2d6+3",
		"Spell: save: 1d20+8, damage: 4d6",
		"Healing: healing: 2d8+4",
		"Multi-attack: to hit: 1d20+7, damage: 1d10+4, to hit: 1d20+7, damage: 1d10+4",
	}
	actionName := "Test Action"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		text := texts[i%len(texts)]
		_, _ = ConvertDiceRolls(text, actionName)
	}
}
