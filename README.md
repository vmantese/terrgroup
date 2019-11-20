# terrgroup
Uses error groups to rapidly iterate through a grouping and asynchronously  
apply an operation to all elements of the group

The result is always thread safe as compared to using an errgroup on a map or struct

Exits early with an error if any of the individual transforms error

Relies on two(or a more optimized third interface)

```go
type Transformer interface{
    Length() (length int)
    Transform(position i) (interface{},error)
}
```

and the grouping that will receive the transformed elements
```go
type Appender interface{
	Append(interface{})
}

```

or the more optimized version that relies on pre-allocated memory where the result cardinality is known
```go
type Injector interface{
	InjectAt(position int, interface{})
}
```


A simple example
```go
package stringhasher


//
//cardinality of result is not known
func hashesJoinedStrings(ss []string) ([]string,error){
	var g terrgroup.Group
	hh := new(HashHolder)
	err := g.Transform(HashesDelimitedStrings(stringsToBeHashed),hh)
	if err != nil{
		return nil,err
	}
	return []string(*hh),nil
} 

//cardinality of result is known
func hashesStrings(ss []string) ([]string,error){
	var g terrgroup.Group
	hashes := make(HashHolder,len(ss))
	err := g.ExactTransform(HashesStrings(stringsToBeHashed),&hashse)
	if err != nil{
		return nil,err
	}
	return hashes,nil
} 

type HashesDelimitedStrings []string

func (h HashesDelimitedStrings) Length() int{
	return len(h)
}

func (h HashesDelimitedStrings) Transform(i int) (interface{},error){
	hasher := sha256.New()
    strings := strings.Split(h[i],",")
    hashes,err := hashesStrings(strings)
    if err != nil{
    	return nil,err
    }
	// successfuly transform
	return hashes,nil
}

type HashesStrings []string

func (h HashesStrings) Length() int{
	return len(h)
}

func (h HashesStrings) Transform(i int) (interface{},error){
	hasher := sha256.New()
	if n,err := encoder.Write([]byte(h[i]));err != nil{
		return nil,err
	}else if n != len(h[i]){
		return nil, errors.New("unexpected number of bytes written to hasher")
	}
	// successfuly transform
	return string(encoder.Sum([]byte{})),nil
}

type HashHolder []string

func (h *HashHolder) Append(i interface{}){
	if str,ok := i.(string);ok{
		*h = append(*h,str)
	}
    if strs,ok := i.([]string);ok{
        *h = append(*h,strs...)
    }
}

func (h HashHolder) InjectAt(i int,j interface{}){
	if str,ok := j.(string);ok{
		h[i] = str
	}
}



```