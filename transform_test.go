package terrgroup

import (
	"crypto/sha256"
	"fmt"
	"reflect"
	"sort"
	"testing"
)






type HashesStrings []string

func (h HashesStrings) Length() int{
	return len(h)
}

func (h HashesStrings) Transform(i int) (interface{},error){
	encoder := sha256.New()
	if _,err := encoder.Write([]byte(h[i]));err != nil{
		return nil,err
	}
	return string(encoder.Sum([]byte{})),nil
}

type HashHolder []string

func (h *HashHolder) Append(i interface{}){
	if str,ok := i.(string);ok{
		*h = append(*h,str)
	}
}

func (h HashHolder) InjectAt(i int,j interface{}){
	if str,ok := j.(string);ok{
		h[i] = str
	}
}

func makeSha256(s string) (string,error){
	encoder := sha256.New()
	if _,err := encoder.Write([]byte(s));err != nil{
		return "",err
	}
	return string(encoder.Sum([]byte{})),nil
}


func tt() []string{
	hh := new(HashHolder)
	return *hh
}

func TestGroup_GoTransform(t *testing.T) {
	var sg Group
	stringsToBeHashed := []string{"first string to be hashed","second string to be hashed"}
	hash1,_ := makeSha256(stringsToBeHashed[0])
	hash2,_ := makeSha256(stringsToBeHashed[1])
	expectedHashes := []string{hash1,hash2}
	actualHashes := HashHolder{}
	err := sg.GoTransform(HashesStrings(stringsToBeHashed),&actualHashes)
	sort.Strings(actualHashes)
	sort.Strings(expectedHashes)
	if err != nil{
		t.Fail()
	}
	if reflect.DeepEqual(expectedHashes,actualHashes){
		fmt.Println(expectedHashes,actualHashes)
		t.Fail()
	}


}

func TestGroup_GoExactTransform(t *testing.T) {
	var sg Group
	stringsToBeHashed := []string{"first string to be hashed","second string to be hashed"}
	hash1,_ := makeSha256(stringsToBeHashed[0])
	hash2,_ := makeSha256(stringsToBeHashed[1])
	expectedHashes := []string{hash1,hash2}
	actualHashes := make(HashHolder,2)
	err := sg.GoExactTransform(HashesStrings(stringsToBeHashed),actualHashes)
	sort.Strings(actualHashes)
	sort.Strings(expectedHashes)
	if err != nil{
		t.Fail()
	}
	if reflect.DeepEqual(expectedHashes,actualHashes){
		fmt.Println(expectedHashes,actualHashes)
		t.Fail()
	}

}


func BenchmarkGroup_GoTransform(b *testing.B) {

}