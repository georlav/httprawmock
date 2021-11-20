package httprawmock

type Route struct {
	Method   string
	Pattern  string
	Response []byte
}

// NewRoute creates a new route object
//
// method can be any valid http method, accepted values are "GET", "HEAD", "POST", "PUT", "PATCH", "DELETE", "CONNECT", "OPTIONS", "TRACE"
// you can also pass an empty method, this will have as result to register the given pattern under all methods
//
// pattern you can use patterns like:
//     /resource/* (match any)
// 	   /resource/{id} (match resource followed by and id)
// 	   /resource/{id:[0-9]+} (match resource followed by a numeric id)
// 	   /resource/{id:[a-z]+} (match resource followed by a alpha id)
// 	   /resource/{id:[a-z0-9]+} (match resource followed by an alphanumeric id)
// 	   /resource/{id}/activate
//     /resource/{id}/resourceb/{uid}
//
// You can use more complex p`atterns, package is using chi router to create routes so if you need more complex patterns you can check their docs
//
// response must a valid raw http response
func NewRoute(method string, pattern string, response []byte) Route {
	return Route{
		Method:   method,
		Pattern:  pattern,
		Response: response,
	}
}
