package cookie

import (
	"net/http"
	"net/url"
	"sync"
)

// Jar implements the http.CookieJar interface. It is used to store cookies for a specific host.
type Jar struct {
	mu  sync.Mutex
	jar map[string][]*http.Cookie
}

// NewJar returns a new Jar.
func NewJar() *Jar {
	jar := &Jar{
		jar: make(map[string][]*http.Cookie),
	}

	return jar
}

// SetCookies ...
func (j *Jar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	j.mu.Lock()
	defer j.mu.Unlock()

	current, ok := j.jar[u.Host]
	if !ok {
		j.jar[u.Host] = cookies
		return
	}

	if compare(current, cookies) {
		return
	}

	j.jar[u.Host] = cookies
}

// Cookies ...
func (j *Jar) Cookies(u *url.URL) []*http.Cookie {
	j.mu.Lock()
	defer j.mu.Unlock()

	return j.jar[u.Host]
}

// AddCookies adds a slice of cookies to the jar.
func (j *Jar) AddCookies(cookies []*http.Cookie) {
	j.mu.Lock()
	defer j.mu.Unlock()

	for _, cookie := range cookies {
		uri, err := url.Parse(cookie.Domain)
		if err != nil {
			continue
		}
		j.jar[cookie.Domain] = append(j.jar[uri.Host], cookie)
	}
}

// AllCookies returns all cookies in the jar.
func (j *Jar) AllCookies() []*http.Cookie {
	j.mu.Lock()
	defer j.mu.Unlock()

	var all []*http.Cookie
	for _, cookies := range j.jar {
		all = append(all, cookies...)
	}

	return all

}

// Clear clears the cookies for a specific host.
func (j *Jar) Clear(u *url.URL) {
	j.mu.Lock()
	defer j.mu.Unlock()

	delete(j.jar, u.Host)
}

// ClearAll clears all cookies in the jar.
func (j *Jar) ClearAll() {
	j.mu.Lock()
	defer j.mu.Unlock()

	j.jar = make(map[string][]*http.Cookie)
}

// compare returns true if the two slices of cookies are equal.
func compare(a []*http.Cookie, b []*http.Cookie) bool {
	for i := range a {
		if a[i].Name != b[i].Name || a[i].Value != b[i].Value {
			return false
		}
	}
	return true
}
