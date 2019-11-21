package examples

import (
	"github.com/vmantese/terrgroup"
	"regexp"
)

//not intended to be comprehensive just illustrative
var SimpleExampleRegex *regexp.Regexp

func init() {
	SimpleExampleRegex = regexp.MustCompile(`\s{0,2}[A-Za-z,;'\\"\s]*((?i)hello|world)[A-Za-z,;'\\"\s]*[.!?]`)
}

type Sentence []byte
type Page []Sentence
type Book []Page

func (p Page) Bytes() []byte {
	byts := make([]byte, 0)
	for i := range p {
		byts = append(byts, p[i]...)
	}
	return byts
}

func ToSentences(bb [][]byte) []Sentence {
	sentences := make([]Sentence, len(bb))
	for i, b := range bb {
		sentences[i] = b
	}
	return sentences
}

//cardinality of result is not known
func ParseBook(book Book) ([]Sentence, error) {
	var g terrgroup.Group
	notepad := new(Notepad)
	if err := g.GoTransform(FindsHelloOrWorld(book), notepad);err != nil{
		return nil,err
	}
	return *notepad,nil
}

//cardinality of result is known
func ParseFirstSentenceBook(book Book) ([]Sentence, error) {
	var g terrgroup.Group
	notepad := make(Notepad, len(book))
	if err := g.GoExactTransform(FindsHelloOrWorld(book), notepad);err != nil{
		return nil,err
	}
	return notepad,nil
}

type FindsHelloOrWorld Book

func (f FindsHelloOrWorld) Length() int {
	//returns number of pages
	return len(f)
}
func (f FindsHelloOrWorld) Transform(i int) (interface{}, error) {
	return findsHelloOrWorld(f[i])
}

func findsHelloOrWorld(page Page) ([]Sentence, error) {
	allBytes := SimpleExampleRegex.FindAll(page.Bytes(), -1)
	return ToSentences(allBytes), nil
}

type Notepad []Sentence

func (n *Notepad) Append(i interface{}) {
	if sentences, ok := i.([]Sentence); ok {
		*n = append(*n, sentences...)
	}
	if sentence, ok := i.(Sentence); ok {
		*n = append(*n, sentence)
	}
}
func (n Notepad) InjectAt(pos int, i interface{}) {
	if sentence, ok := i.(Sentence); ok {
		n[pos] = sentence
	}
	//just take the first
	if sentences, ok := i.([]Sentence); ok && len(sentences) > 0 {
		n[pos] = sentences[0]
	}
}
