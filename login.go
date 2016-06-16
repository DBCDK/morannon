package main

import (
	"html/template"
	"net/http"
	"time"
)

func showLogin(w http.ResponseWriter, r *http.Request) {
	type LoginData struct{}

	var page = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>Marathon Login</title>
	</head>
	<body>
		<form method="POST" action="/login">
			<input type="text" name="access_token">
			<input type="submit" value="login">
		</form>
	</body>
</html>`

	t, err := template.New("t").Parse(page)
	if err != nil {
		panic(err)
	}

	data := LoginData{}

	t.Execute(w, data)
}

func performLogin(w http.ResponseWriter, r *http.Request) {
	token := r.PostFormValue("access_token")

	expiration := time.Now().Add(7 * 24 * time.Hour)

	cookie := http.Cookie{Name: "access_token", Value: token, Expires: expiration}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/", 302)
}
