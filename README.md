# GO-CHECK-TOOL 
go-check-tool is a library that allows you to check fields of type `string` in a structure.

## Guide
### Installation
```bach
    go get github.com/samuskitchen/go-check-tool
```

Example
```go 
    type Person struct {
        IdentityDocument    string `chk:"num len=40"`
        FirstName           string `chk:"max=255"`
        LastName            string `chk:"min=10 max=255"`
        Age                 uint
    }

    func main() {
        p := Person{IdentityDocument: "98736712", FirstName: "Kevin", LastName: "Saucedo", Age: 23}
        //return err if `p` is invalid
        if err := check_tool.Valid(p); err != nil {
            log.Println(err)
        }
        
        //valid with omits, Field FirstName is omited
        if err := check_tool.ValidWithOmit(p, kcheck.OmitFields{"FirstName"}); err != nil {
            log.Println(err)
        }
    }
```

| tag-keys          | Description                                                                                                                                                                                                                                                            |
|-------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `nonil`           | the string cannot be null or empty, and empty spaces don't count                                                                                                                                                                                                       |
| `nosp`            | checks that there are no empty spaces at the beginning or end of the string                                                                                                                                                                                            |
| `word`            | values `0-9` `a-z` `A-Z` and `_` are the only ones allowed, (no spaces)                                                                                                                                                                                                |
| `txt`             | text, safe text; only allows text with no leading or trailing spaces, words cannot be separated by more than two spaces, and only characters not in this list are allowed: `!\"#$%&'()*+, ./:;<=>?@[\\]^_}{~\`<br/> can be used to receive name, surname among others  |
| `email`           | email                                                                                                                                                                                                                                                                  |
| `url`             | Web URL, HTTP, and HTTPS (not yet available)                                                                                                                                                                                                                           |
| `num`             | check if all characters are numeric                                                                                                                                                                                                                                    |
| `decimal`         | only decimal numbers; `00.00`, `1.00` is valid                                                                                                                                                                                                                         |
| `len=:number`     | the length of the entered string must be equal to `:number`                                                                                                                                                                                                            |
| `max=:number`     | max length must be less than `:number`                                                                                                                                                                                                                                 |
| `min=:number`     | minimum length must be greater than `:number`                                                                                                                                                                                                                          |
| `rgx=:expression` | allows passing a regular expression, e.g. key-tag: `regex=(^[0-9]{8}$)`                                                                                                                                                                                                |
--------------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------------
Tag example: `chk:"nonil num len=8"` `the string cannot be empty, it can only be 8 numeric characters`