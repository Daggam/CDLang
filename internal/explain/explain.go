package explain

import (
	"fmt"
	"strings"

	"github.com/Daggam/CDLang/internal/ast"
	"github.com/Daggam/CDLang/internal/lexer"
	"github.com/Daggam/CDLang/internal/parser"
	"github.com/Daggam/CDLang/internal/token"
)

func Build(input string) string {
	var out strings.Builder

	out.WriteString("[1] CODIGO FUENTE\n")
	out.WriteString(strings.TrimSpace(input))
	out.WriteString("\n\n")

	out.WriteString("[2] TOKENS\n")
	for _, tok := range collectTokens(input) {
		out.WriteString(formatToken(tok))
		out.WriteString("\n")
	}
	out.WriteString("\n")

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) > 0 {
		out.WriteString("[3] AST\n")
		out.WriteString("No se pudo construir el AST.\n\n")
		out.WriteString("ERRORES DEL PARSER\n")
		for _, err := range p.Errors() {
			out.WriteString("- ")
			out.WriteString(err)
			out.WriteString("\n")
		}
		return out.String()
	}

	explainStmt := firstExplainStatement(program)
	if explainStmt == nil {
		out.WriteString("[3] AST\n")
		out.WriteString("El comando no es un EXPLAIN valido.\n")
		return out.String()
	}

	out.WriteString("[3] AST\n")
	out.WriteString(formatExplainAST(explainStmt))
	out.WriteString("\n")

	out.WriteString("[4] VALIDACION SEMANTICA\n")
	out.WriteString(formatSemanticValidation(explainStmt.Inner))
	out.WriteString("\n")

	out.WriteString("[5] PLAN DE EJECUCION\n")
	out.WriteString(formatExecutionPlan(explainStmt.Inner))
	out.WriteString("\n")

	out.WriteString("[6] RESULTADO\n")
	out.WriteString(formatResult(explainStmt.Inner))

	return out.String()
}

func IsExplain(input string) bool {
	l := lexer.New(input)
	return l.NextToken().Type == token.EXPLAIN
}

func collectTokens(input string) []token.Token {
	l := lexer.New(input)
	tokens := []token.Token{}
	for {
		tok := l.NextToken()
		tokens = append(tokens, tok)
		if tok.Type == token.EOF {
			break
		}
	}
	return tokens
}

func formatToken(tok token.Token) string {
	switch tok.Type {
	case token.IDENT, token.INT:
		return fmt.Sprintf("%s(%s)", tok.Type, tok.Literal)
	case token.SEMICOLON:
		return "SEMICOLON"
	case token.COMMA:
		return "COMMA"
	case token.LPAREN:
		return "LPAREN"
	case token.RPAREN:
		return "RPAREN"
	case token.ASSIGN:
		return "ASSIGN"
	case token.PLUS:
		return "PLUS"
	case token.EOF:
		return "EOF"
	default:
		return string(tok.Type)
	}
}

func firstExplainStatement(program *ast.Program) *ast.ExplainStatement {
	if program == nil || len(program.Statements) == 0 {
		return nil
	}
	stmt, ok := program.Statements[0].(*ast.ExplainStatement)
	if !ok {
		return nil
	}
	return stmt
}

func formatExplainAST(stmt *ast.ExplainStatement) string {
	var out strings.Builder
	out.WriteString("ExplainStatement {\n")
	out.WriteString("  Inner: ")
	out.WriteString(formatStatementAST(stmt.Inner, "  "))
	out.WriteString("}\n")
	return out.String()
}

func formatStatementAST(stmt ast.Statement, indent string) string {
	switch s := stmt.(type) {
	case *ast.OfferStatement:
		return fmt.Sprintf("OfferStatement {\n%s  Collectables: %s\n%s}\n", indent, formatCollectables(s.Collectables), indent)
	case *ast.GetOfferStatement:
		return fmt.Sprintf("GetOfferStatement {\n%s  Identifier: %s\n%s}\n", indent, s.Identifier.Value, indent)
	case *ast.SendOfferStatement:
		return fmt.Sprintf("SendOfferStatement {\n%s  Offered: %s\n%s  Requested: %s\n%s  User: %s\n%s}\n",
			indent, formatCollectables(s.LCollectables), indent, formatCollectables(s.RCollectables), indent, s.Username.Value, indent)
	case *ast.ViewOfferStatement:
		return fmt.Sprintf("ViewOfferStatement {\n%s  Target: OFFER\n%s}\n", indent, indent)
	case *ast.AcceptTradeStatement:
		return fmt.Sprintf("AcceptTradeStatement {\n%s  OfferID: %s\n%s}\n", indent, formatIDs(s.OfferID), indent)
	case *ast.DeclineTradeStatement:
		return fmt.Sprintf("DeclineTradeStatement {\n%s  OfferID: %s\n%s}\n", indent, formatIDs(s.OfferID), indent)
	case *ast.DeleteOfferStatement:
		return fmt.Sprintf("DeleteOfferStatement {\n%s  Collectables: %s\n%s}\n", indent, formatCollectables(s.Collectables), indent)
	default:
		return fmt.Sprintf("%T\n", stmt)
	}
}

func formatSemanticValidation(stmt ast.Statement) string {
	switch s := stmt.(type) {
	case *ast.OfferStatement:
		return strings.Join([]string{
			"- Verificar que el usuario actual exista.",
			"- Verificar que el usuario posea: " + formatCollectables(s.Collectables) + ".",
			"- Verificar que los codigos de coleccionables existan en el catalogo.",
			"- Verificar que las cantidades sean mayores que cero.",
		}, "\n") + "\n"
	case *ast.GetOfferStatement:
		return strings.Join([]string{
			"- Verificar que el coleccionable " + s.Identifier.Value + " exista.",
			"- Buscar coleccionables intercambiables con ese nombre.",
			"- Excluir coleccionables del usuario actual.",
		}, "\n") + "\n"
	case *ast.SendOfferStatement:
		return strings.Join([]string{
			"- Verificar que el usuario destino exista: " + s.Username.Value + ".",
			"- Verificar que el usuario actual posea lo ofrecido: " + formatCollectables(s.LCollectables) + ".",
			"- Verificar que el usuario destino tenga disponible lo solicitado: " + formatCollectables(s.RCollectables) + ".",
		}, "\n") + "\n"
	case *ast.ViewOfferStatement:
		return strings.Join([]string{
			"- Verificar que el usuario actual exista.",
			"- Buscar ofertas o coleccionables intercambiables visibles.",
		}, "\n") + "\n"
	case *ast.AcceptTradeStatement:
		return validationForTrade("aceptar", s.OfferID)
	case *ast.DeclineTradeStatement:
		return validationForTrade("rechazar", s.OfferID)
	case *ast.DeleteOfferStatement:
		return strings.Join([]string{
			"- Verificar que el usuario actual exista.",
			"- Verificar que existan ofertas propias con: " + formatCollectables(s.Collectables) + ".",
			"- Verificar que esas ofertas sigan activas.",
		}, "\n") + "\n"
	default:
		return "- No hay validacion semantica definida para este comando.\n"
	}
}

func validationForTrade(action string, ids []int) string {
	target := "los trades indicados"
	if len(ids) == 0 {
		target = "todos los trades pendientes aplicables"
	}
	return strings.Join([]string{
		"- Verificar que existan " + target + ": " + formatIDs(ids) + ".",
		"- Verificar que cada trade este en estado pending.",
		"- Verificar permisos del usuario actual para " + action + ".",
		"- Verificar disponibilidad de los coleccionables involucrados.",
	}, "\n") + "\n"
}

func formatExecutionPlan(stmt ast.Statement) string {
	switch s := stmt.(type) {
	case *ast.OfferStatement:
		return numbered([]string{
			"Tomar los coleccionables: " + formatCollectables(s.Collectables) + ".",
			"Validar stock del usuario actual.",
			"Marcar esos coleccionables como intercambiables.",
		})
	case *ast.GetOfferStatement:
		return numbered([]string{
			"Buscar coleccionables intercambiables para " + s.Identifier.Value + ".",
			"Mostrar la tabla de resultados.",
		})
	case *ast.SendOfferStatement:
		return numbered([]string{
			"Buscar al usuario destino " + s.Username.Value + ".",
			"Validar lo ofrecido: " + formatCollectables(s.LCollectables) + ".",
			"Validar lo solicitado: " + formatCollectables(s.RCollectables) + ".",
			"Crear una propuesta de trade.",
		})
	case *ast.ViewOfferStatement:
		return numbered([]string{
			"Buscar ofertas o coleccionables intercambiables.",
			"Mostrar los resultados al usuario.",
		})
	case *ast.AcceptTradeStatement:
		return tradePlan("aceptar", "accepted", s.OfferID)
	case *ast.DeclineTradeStatement:
		return tradePlan("rechazar", "declined", s.OfferID)
	case *ast.DeleteOfferStatement:
		return numbered([]string{
			"Buscar ofertas propias que contengan: " + formatCollectables(s.Collectables) + ".",
			"Eliminar o desactivar esas ofertas.",
		})
	default:
		return "1. No hay plan definido para este comando.\n"
	}
}

func tradePlan(action string, finalStatus string, ids []int) string {
	return numbered([]string{
		"Buscar " + formatIDs(ids) + ".",
		"Validar permisos del usuario actual para " + action + ".",
		"Validar estado y disponibilidad de coleccionables.",
		"Cambiar el trade al estado " + finalStatus + ".",
	})
}

func formatResult(stmt ast.Statement) string {
	switch s := stmt.(type) {
	case *ast.OfferStatement:
		return "Se marcarian como intercambiables: " + formatCollectables(s.Collectables) + ".\n"
	case *ast.GetOfferStatement:
		return "Se mostraria la disponibilidad de intercambios para: " + s.Identifier.Value + ".\n"
	case *ast.SendOfferStatement:
		return "Se crearia una propuesta para " + s.Username.Value + ": " + formatCollectables(s.LCollectables) + " por " + formatCollectables(s.RCollectables) + ".\n"
	case *ast.ViewOfferStatement:
		return "Se mostraria el listado visible para el usuario actual.\n"
	case *ast.AcceptTradeStatement:
		return "Se aceptarian los trades: " + formatIDs(s.OfferID) + ".\n"
	case *ast.DeclineTradeStatement:
		return "Se rechazarian los trades: " + formatIDs(s.OfferID) + ".\n"
	case *ast.DeleteOfferStatement:
		return "Se eliminarian ofertas propias asociadas a: " + formatCollectables(s.Collectables) + ".\n"
	default:
		return "No hay resultado definido para este comando.\n"
	}
}

func numbered(steps []string) string {
	var out strings.Builder
	for i, step := range steps {
		fmt.Fprintf(&out, "%d. %s\n", i+1, step)
	}
	return out.String()
}

func formatCollectables(collectables []*ast.Collectable) string {
	parts := make([]string, 0, len(collectables))
	for _, collectable := range collectables {
		parts = append(parts, fmt.Sprintf("%s(%d)", collectable.Value, collectable.Amount))
	}
	return strings.Join(parts, ", ")
}

func formatIDs(ids []int) string {
	if len(ids) == 0 {
		return "sin ID especifico"
	}

	parts := make([]string, 0, len(ids))
	for _, id := range ids {
		parts = append(parts, fmt.Sprintf("%d", id))
	}
	return strings.Join(parts, ", ")
}
