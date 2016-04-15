# beego-response
This is an additional golang package for the [beego web framework](https://github.com/astaxie/beego) to simplify error and response handling within controllers. It wraps functionality of [*BeegoOutput](https://github.com/astaxie/beego/blob/master/context/output.go).

**Work in progress!**
Package is used as an example in https://github.com/moehlone/mongodm-example

[![GoDoc](https://godoc.org/github.com/zebresel-com/beego-response?status.svg)](https://godoc.org/github.com/zebresel-com/beego-response)

### Advantages
- pagination
- direct JSON error response
- custom data appending

### Example
```go
  func (self *UsersController) Me() {
	  if self.token != nil && self.user != nil {
		  self.response.AddContent("user", self.user)
		  self.response.ServeJSON()
	  } else {
	    self.response.Error(http.StatusBadRequest)
	  }
  }
```
**Feel free to contribute!**
