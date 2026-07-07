package explain

import (
	"strings"
	"testing"
)

func TestBuildExplainOutput(t *testing.T) {
	output := Build("EXPLAIN OFFER AR-LM10;")

	expectedParts := []string{
		"[1] CODIGO FUENTE",
		"[2] TOKENS",
		"EXPLAIN",
		"IDENT(AR-LM10)",
		"[3] AST",
		"ExplainStatement",
		"OfferStatement",
		"[4] VALIDACION SEMANTICA",
		"[5] PLAN DE EJECUCION",
		"[6] RESULTADO",
	}

	for _, part := range expectedParts {
		if !strings.Contains(output, part) {
			t.Fatalf("La salida de EXPLAIN no contiene %q.\nSalida:\n%s", part, output)
		}
	}
}

func TestIsExplain(t *testing.T) {
	if !IsExplain("EXPLAIN VIEW OFFER;") {
		t.Fatal("Se esperaba que IsExplain detecte EXPLAIN.")
	}

	if IsExplain("VIEW OFFER;") {
		t.Fatal("No se esperaba que IsExplain detecte EXPLAIN.")
	}
}
