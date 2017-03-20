package rdf

import (
	"reflect"
	"testing"
)

func TestGraphIsomorphism(t *testing.T) {
	tests := []struct {
		a  string
		b  string
		eq bool
	}{
		{
			"",
			"",
			true,
		},
		{
			`<s> <p> "o" .`,
			`<s> <p> "o" .`,
			true,
		},
		{
			`<s> <p> "o" .`,
			`<s> <p> "o"^^<http://www.w3.org/2001/XMLSchema#string> .`,
			true,
		},
		{
			`<s> <p> "o" .`,
			`<s> <p> "o2" .`,
			false,
		},
		{
			`<s> <p> "o" .
			 <s> <p> "xyz" .`,
			`<s> <p> "o" .`,
			false,
		},
		{
			`_:a <p> "1" .`,
			`_:a <p> "1" .`,
			true,
		},
		{
			`_:a <p> "1" .`,
			`_:b <p> "1" .`,
			true,
		},
		{
			`<s> <p> _:a .`,
			`<s> <p> _:a .`,
			true,
		},
		{
			`<s> <p> _:a .`,
			`<s> <p> _:b .`,
			true,
		},
		{
			`<s> <p> _:a .
			 _:a <p> "xyz" . `,
			`<s> <p> _:b .
			 _:b <p> "xyz" . `,
			true,
		},
		{
			`<s> <p> _:a .
			 _:a <p> "xyz" .
			 _:a <p> "æøå" . `,
			`<s> <p> _:a .
			_:a <p> "xyz" . `,
			false,
		},
		{
			`<s> <p> _:a .
			 _:a <p> "xyz" .
			 _:a <p> "æøå" .`,
			`<s> <p> _:b .
			 _:b <p> "xyz" .
			 _:b <p> "æøå" .`,
			true,
		},
		{
			`<s> <p> _:a .
			 _:a <p> "xyz" . `,
			`<s> <p> _:a .
			 _:a <p> "XYZ" . `,
			false,
		},
		{
			`<s> <p> _:a .
			 _:a <p> "xyz" .
			 _:a <p> "æøå" .`,
			`<s> <p> _:b .
			 _:b <p> "xyz" .
			 _:b <p> "æøå!!!" .`,
			false,
		},
		{
			`<s> <p> _:a .
			 _:a <p> "xyz" .
			 <s> <p> _:a2 .
			 _:a2 <p> "123" .`,
			`<s> <p> _:a .
			 _:a <p> "xyz" .
			 <s> <p> _:a2 .
			 _:a2 <p> "123" .`,
			true,
		},
	}

	for _, test := range tests {
		a := mustDecode(test.a)
		b := mustDecode(test.b)
		if got := a.Eq(b); got != test.eq {
			t.Fatalf("\n%v  Eq\n%v  = %v; want %v", mustEncode(a), mustEncode(b), got, test.eq)
		}
	}
}

func TestGroupPatternsByVariable(t *testing.T) {
	tests := []struct {
		patterns []TriplePattern
		want     [][]TriplePattern
	}{
		{
			[]TriplePattern{},
			[][]TriplePattern{},
		},
		{
			[]TriplePattern{
				{NewVariable("s"), NewVariable("p"), NewVariable("o")},
			},
			[][]TriplePattern{{
				{NewVariable("s"), NewVariable("p"), NewVariable("o")},
			}},
		},
		{
			[]TriplePattern{
				{mustNamedNode("h1"), mustNamedNode("knows"), NewVariable("a")},
				{NewVariable("b"), mustNamedNode("knows"), mustNamedNode("h2")},
			},
			[][]TriplePattern{
				{{mustNamedNode("h1"), mustNamedNode("knows"), NewVariable("a")}},
				{{NewVariable("b"), mustNamedNode("knows"), mustNamedNode("h2")}}},
		},
		{
			[]TriplePattern{
				{mustNamedNode("h1"), mustNamedNode("knows"), NewVariable("a")},
				{NewVariable("a"), mustNamedNode("knows"), mustNamedNode("h2")},
			},
			[][]TriplePattern{{
				{mustNamedNode("h1"), mustNamedNode("knows"), NewVariable("a")},
				{NewVariable("a"), mustNamedNode("knows"), mustNamedNode("h2")},
			}},
		},
		{
			[]TriplePattern{
				{mustNamedNode("h1"), mustNamedNode("knows"), NewVariable("a")},
				{NewVariable("a"), mustNamedNode("knows"), NewVariable("b")},
				{NewVariable("b"), mustNamedNode("knows"), mustNamedNode("h2")},
				{NewVariable("c"), mustNamedNode("related"), mustNamedNode("h1")},
			},
			[][]TriplePattern{{{mustNamedNode("h1"), mustNamedNode("knows"), NewVariable("a")},
				{NewVariable("a"), mustNamedNode("knows"), NewVariable("b")},
				{NewVariable("b"), mustNamedNode("knows"), mustNamedNode("h2")},
			},
				{{NewVariable("c"), mustNamedNode("related"), mustNamedNode("h1")}}},
		},
		{
			[]TriplePattern{
				{mustNamedNode("h1"), mustNamedNode("knows"), NewVariable("a")},
				{NewVariable("a"), mustNamedNode("knows"), NewVariable("b")},
				{NewVariable("b"), mustNamedNode("knows"), mustNamedNode("h2")},
				{NewVariable("c"), mustNamedNode("related"), NewVariable("a")},
			},
			[][]TriplePattern{{{mustNamedNode("h1"), mustNamedNode("knows"), NewVariable("a")},
				{NewVariable("a"), mustNamedNode("knows"), NewVariable("b")},
				{NewVariable("b"), mustNamedNode("knows"), mustNamedNode("h2")},
				{NewVariable("c"), mustNamedNode("related"), NewVariable("a")}},
			},
		},
	}

	for _, test := range tests {
		if groups := groupPatternsByVariable(test.patterns); !reflect.DeepEqual(groups, test.want) {
			t.Errorf("groupPatternsByVariable(%v) => %v; want %v", test.patterns, groups, test.want)
		}
	}
}

func TestGraphWhere(t *testing.T) {
	g := mustDecode(`
		<a1> <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <Person> .
		<a1> <hasName> "Italo Calvino" .
		<a1> <hasBirthYear> "1923"^^<http://www.w3.org/2001/XMLSchema#gYear> .
		<a1> <hasDeathYear> "1985"^^<http://www.w3.org/2001/XMLSchema#gYear> .

		<w1> <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <Work> .
		<w1> <hasMainTitle> "Le Cosmicomiche" .
		<w1> <hasPublicationYear> "1965"^^<http://www.w3.org/2001/XMLSchema#gYear> .
		<w1> <hasContributor> _:c1 .

		_:c1 <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <Contribution> .
		_:c1 <hasRole> <author> .
		_:c1 <hasAgent> <a1> .

		<p1> <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <Publication> .
		<p1> <isPublicationOf> <w1> .
		<p1> <hasMainTitle> "The Complete Cosmicomics" .
		<p1> <hasISBN> "9781846141652" .
		<p1> <hasPublishYear> "2009"^^<http://www.w3.org/2001/XMLSchema#gYear> .
		<p1> <isPublishedBy> <c1> .
		<p1> <hasContributor> _:c2 .
		<p1> <hasContributor> _:c3 .
		<p1> <hasContributor> _:c4 .
		<p1> <hasAbstract> "The definitive edition of Calvino’s cosmicomics, bringing together all of these enchanting stories—including some never before translated — in one volume for the first time" .

		_:c2 <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <Contribution> .
		_:c2 <hasRole> <translator> .
		_:c2 <hasAgent> <a2> .

		_:c3 <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <Contribution> .
		_:c3 <hasRole> <translator> .
		_:c3 <hasAgent> <a3> .

		_:c4 <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <Contribution> .
		_:c4 <hasRole> <translator> .
		_:c4 <hasAgent> <a4> .

		<a2> <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <Person> .
		<a2> <hasName> "Martin L. McLaughlin" .

		<a3> <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <Person> .
		<a3> <hasName> "Tim Parks" .

		<a4> <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <Person> .
		<a4> <hasName> "William Weaver" .

		<c1> <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <Corporation> .
		<c1> <hasName> "Penguin" .

		<w2> <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <Work> .
		<w2> <hasMainTitle> "Il barone rampante" .
		<w2> <hasPublicationYear> "1957"^^<http://www.w3.org/2001/XMLSchema#gYear> .
		<w2> <hasContributor> _:c5 .

		_:c5 <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <Contribution> .
		_:c5 <hasRole> <author> .
		_:c5 <hasAgent> <a1> .

		<p2> <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <Publication> .
		<p2> <isPublicationOf> <w2> .
		<p2> <hasMainTitle> "Klatrebaronen" .
		<p2> <hasPublishYear> "1961"^^<http://www.w3.org/2001/XMLSchema#gYear> .
		<p2> <isPublishedBy> <c2> .
		<p2> <hasContributor> _:c6 .

		_:c6 <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <Contribution> .
		_:c6 <hasRole> <translator> .
		_:c6 <hasAgent> <a4> .

		<a5> <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <Person> .
		<a5> <hasName> "Ingeborg Hagemann" .

		<c1> <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <Corporation> .
		<c1> <hasName> "Gyldendal" .
		`)

	tests := []struct {
		patterns []TriplePattern
		want     string
	}{
		{
			[]TriplePattern{},
			"",
		},
		{
			[]TriplePattern{
				{NewVariable("s"), NewVariable("p"), NewVariable("o")},
			},
			mustEncode(g),
		},
		{
			[]TriplePattern{
				{
					NewVariable("pub"),
					mustNamedNode("hasMainTitle"),
					NewVariable("title"),
				},
			},
			`
			<w1> <hasMainTitle> "Le Cosmicomiche" .
			<p1> <hasMainTitle> "The Complete Cosmicomics" .
			<w2> <hasMainTitle> "Il barone rampante" .
			<p2> <hasMainTitle> "Klatrebaronen" .`,
		},
		{
			[]TriplePattern{
				{
					mustNamedNode("p2"),
					mustNamedNode("hasPublishYear"),
					NewTypedLiteral("1961", mustNamedNode("http://www.w3.org/2001/XMLSchema#gYear")),
				},
			},
			`<p2> <hasPublishYear> "1961"^^<http://www.w3.org/2001/XMLSchema#gYear> .`,
		},
		{
			[]TriplePattern{
				{
					mustNamedNode("p2"),
					NewVariable("p"),
					NewVariable("o"),
				},
			},
			`
			<p2> <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <Publication> .
			<p2> <isPublicationOf> <w2> .
			<p2> <hasMainTitle> "Klatrebaronen" .
			<p2> <hasPublishYear> "1961"^^<http://www.w3.org/2001/XMLSchema#gYear> .
			<p2> <isPublishedBy> <c2> .
			<p2> <hasContributor> _:1 .`,
		},
		{
			[]TriplePattern{
				{
					NewVariable("s"),
					NewVariable("p"),
					NewStrLiteral("Penguin"),
				},
			},
			`<c1> <hasName> "Penguin" .`,
		},
		{
			[]TriplePattern{
				{
					NewVariable("s"),
					NewVariable("p"),
					NewLangLiteral("Penguin", "en"),
				},
			},
			"",
		},
		{
			[]TriplePattern{
				{
					NewVariable("w"),
					mustNamedNode("hasMainTitle"),
					NewStrLiteral("Le Cosmicomiche"),
				},
				{
					NewVariable("p"),
					mustNamedNode("isPublicationOf"),
					NewVariable("w"),
				},
				{
					NewVariable("p"),
					mustNamedNode("hasMainTitle"),
					NewVariable("pubTitle"),
				},
			},
			`
			<w1> <hasMainTitle> "Le Cosmicomiche" .
			<p1> <isPublicationOf> <w1> .
			<p1> <hasMainTitle> "The Complete Cosmicomics" .
			`,
		},
	}

	for _, test := range tests {
		if got := g.Where(test.patterns); !got.Eq(mustDecode(test.want)) {
			t.Errorf("got:\n%v\nwant:\n%v", mustEncode(got), test.want)
		}
	}
}

func mustNamedNode(s string) NamedNode {
	node, err := NewNamedNode(s)
	if err != nil {
		panic(err)
	}
	return node
}
