package main

import (
    "fmt"
    "io"
    "log"
    "net/http"
)

// const (
//     version         = "v1.0.0"
//     authServiceURL  = "http://localhost:8081/auth"
//     colorServiceURL = "http://localhost:8082/color"
// )

const (
    version         = "v2.0.0"
    authServiceURL  = "http://auth.zeinfahrozi.my.id/auth"
    colorServiceURL = "http://color.zeinfahrozi.my.id/color"
)

func main() {
    http.HandleFunc("/", serveHome)
    http.HandleFunc("/login", serveLogin)
    http.HandleFunc("/do_login", handleLogin)
    http.HandleFunc("/protected", serveProtected)
    http.HandleFunc("/get_color", handleGetColor)

    log.Printf("HTMX Frontend (version %s) running on :8080\n", version)
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func serveHome(w http.ResponseWriter, r *http.Request) {
    html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <title>HTMX Frontend (Created By Zein)</title>
    <script src="https://unpkg.com/htmx.org@1.9.2"></script>
</head>
<body>
    <h1>Welcome</h1>
    <p><em>Version: %s</em></p>
    <a href="/login">Login to access protected page</a>
</body>
</html>`, version)
    w.Header().Set("Content-Type", "text/html")
    fmt.Fprint(w, html)
}

func serveLogin(w http.ResponseWriter, r *http.Request) {
    html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <title>Login</title>
    <script src="https://unpkg.com/htmx.org@1.9.2"></script>
</head>
<body>
    <h2>Login</h2>
    <p><em>Version: %s</em></p>
    <form hx-post="/do_login" hx-target="#login_result" hx-swap="innerHTML">
        <input type="text" name="username" placeholder="Username" required><br>
        <input type="password" name="password" placeholder="Password" required><br>
        <button type="submit">Login</button>
    </form>
    <div id="login_result"></div>
</body>
</html>`, version)
    w.Header().Set("Content-Type", "text/html")
    fmt.Fprint(w, html)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
    username := r.FormValue("username")
    password := r.FormValue("password")

    req, err := http.NewRequest("POST", authServiceURL, nil)
    if err != nil {
        http.Error(w, "Internal error", http.StatusInternalServerError)
        return
    }
    req.SetBasicAuth(username, password)
    resp, err := http.DefaultClient.Do(req)
    if err != nil || resp.StatusCode != http.StatusOK {
        fmt.Fprint(w, "<span style='color:red'>Login failed</span>")
        return
    }

    http.SetCookie(w, &http.Cookie{
        Name:  "session_user",
        Value: username,
        Path:  "/",
    })
    fmt.Fprintf(w, `<span style="color:green">Login successful! <a href="/protected">Go to protected page</a></span><br><em>Version: %s</em>`, version)
}

func serveProtected(w http.ResponseWriter, r *http.Request) {
    cookie, err := r.Cookie("session_user")
    if err != nil || cookie.Value == "" {
        http.Redirect(w, r, "/login", http.StatusFound)
        return
    }

    html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <title>Protected Page</title>
    <script src="https://unpkg.com/htmx.org@1.9.2"></script>
</head>
<body>
    <h2>Protected Page</h2>
    <p><em>Version: %s</em></p>
    <p>Welcome, %s!</p>
    <button hx-get="/get_color" hx-target="#color_result" hx-swap="innerHTML">Get Random Color</button>
    <div id="color_result"></div>
    <br>
    <a href="/">Back to Home</a>
</body>
</html>`, version, cookie.Value)
    w.Header().Set("Content-Type", "text/html")
    fmt.Fprint(w, html)
}

func handleGetColor(w http.ResponseWriter, r *http.Request) {
    resp, err := http.Get(colorServiceURL)
    if err != nil {
        http.Error(w, "Failed to get color", http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()
    color, _ := io.ReadAll(resp.Body)
    fmt.Fprintf(w, "<div style='color:%s'>Random Color: %s<br><em>Version: %s</em></div>", color, color, version)
}
