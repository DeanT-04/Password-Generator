package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"html/template"
	"math/big"
	"net/http"
	"strconv"
)

const htmlTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Password Generator</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            background-color: #0f0f0f;
            color: #fff;
            font-family: monospace;
            min-height: 100vh;
            display: flex;
            justify-content: center;
            align-items: center;
            padding: 2rem 0;
        }

        .container {
            background-color: #1a1a1a;
            padding: 2rem;
            border-radius: 10px;
            border: 2px solid #ffeb3b;
            box-shadow: 0 0 20px rgba(255, 235, 59, 0.2);
            width: 90%;
            max-width: 500px;
        }

        h1 {
            text-align: center;
            color: #ffeb3b;
            margin-bottom: 2rem;
            text-shadow: 0 0 10px rgba(255, 235, 59, 0.5);
        }

        #password-box {
            background-color: #000;
            border: 1px solid #4caf50;
            padding: 1rem;
            margin: 1rem 0;
            border-radius: 4px;
            color: #4caf50;
            cursor: pointer;
            font-family: monospace;
            text-align: center;
            min-height: 3rem;
            display: flex;
            align-items: center;
            justify-content: center;
        }

        #password-box:hover {
            background-color: #0a0a0a;
        }

        .controls {
            display: grid;
            gap: 1rem;
            margin: 1rem 0;
        }

        .control-group {
            display: flex;
            align-items: center;
            gap: 0.5rem;
        }

        input[type="number"] {
            background-color: #000;
            border: 1px solid #ffeb3b;
            color: #fff;
            padding: 0.5rem;
            border-radius: 4px;
            width: 5rem;
        }

        button {
            background-color: #ffeb3b;
            color: #000;
            border: none;
            padding: 0.75rem 1.5rem;
            border-radius: 4px;
            cursor: pointer;
            font-weight: bold;
            width: 100%;
            margin-top: 1rem;
        }

        button:hover {
            background-color: #ffd700;
        }

        .password-history {
            margin-top: 2rem;
            border-top: 1px solid #333;
            padding-top: 1rem;
        }

        .password-history h2 {
            color: #ffeb3b;
            font-size: 1.2rem;
            margin-bottom: 1rem;
        }

        #history-list {
            max-height: 120px;
            overflow-y: auto;
            padding-right: 5px;
        }

        #history-list::-webkit-scrollbar {
            width: 8px;
        }

        #history-list::-webkit-scrollbar-track {
            background: #1a1a1a;
            border-radius: 4px;
        }

        #history-list::-webkit-scrollbar-thumb {
            background: #4caf50;
            border-radius: 4px;
        }

        #history-list::-webkit-scrollbar-thumb:hover {
            background: #45a049;
        }

        .history-item {
            background-color: #000;
            border: 1px solid #4caf50;
            padding: 0.5rem;
            margin-bottom: 0.5rem;
            border-radius: 4px;
            color: #4caf50;
            font-family: monospace;
            height: 36px;
            display: flex;
            align-items: center;
            justify-content: space-between;
        }

        .history-item .password-text {
            flex-grow: 1;
            margin-right: 10px;
            cursor: pointer;
        }

        .delete-btn {
            color: #ff0000;
            font-weight: bold;
            cursor: pointer;
            font-family: Arial, sans-serif;
            font-size: 18px;
            width: 24px;
            height: 24px;
            display: flex;
            align-items: center;
            justify-content: center;
            border-radius: 50%;
            transition: all 0.2s ease;
        }

        .delete-btn:hover {
            background-color: #ff0000;
            color: #fff;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>Password Generator</h1>
        <div id="password-box" onclick="copyToClipboard(this)"></div>
        <form id="generator-form">
            <div class="controls">
                <div class="control-group">
                    <label for="length">Length:</label>
                    <input type="number" id="length" name="length" value="12" min="4" max="128">
                </div>
                <div class="control-group">
                    <input type="checkbox" id="uppercase" name="uppercase" checked>
                    <label for="uppercase">Uppercase Letters</label>
                </div>
                <div class="control-group">
                    <input type="checkbox" id="lowercase" name="lowercase" checked>
                    <label for="lowercase">Lowercase Letters</label>
                </div>
                <div class="control-group">
                    <input type="checkbox" id="numbers" name="numbers" checked>
                    <label for="numbers">Numbers</label>
                </div>
                <div class="control-group">
                    <input type="checkbox" id="special" name="special" checked>
                    <label for="special">Special Characters</label>
                </div>
            </div>
            <button type="submit">Generate Password</button>
        </form>
        <div class="password-history">
            <h2>Password History:</h2>
            <div id="history-list">
                {{range .History}}
                <div class="history-item">
                    <span class="password-text" onclick="copyToClipboard(this)">{{.}}</span>
                    <span class="delete-btn" onclick="deletePassword(this)">&#215;</span>
                </div>
                {{end}}
            </div>
        </div>
    </div>
    <script>
        function copyToClipboard(element) {
            const text = element.textContent;
            navigator.clipboard.writeText(text);
            
            const originalColor = element.style.color;
            element.style.color = '#fff';
            setTimeout(() => {
                element.style.color = originalColor;
            }, 200);
        }

        function deletePassword(element) {
            const historyItem = element.parentElement;
            const password = historyItem.querySelector('.password-text').textContent;
            
            fetch('/delete', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ password: password })
            })
            .then(response => {
                if (response.ok) {
                    historyItem.remove();
                }
            });
        }

        document.getElementById('generator-form').onsubmit = function(e) {
            e.preventDefault();
            const form = e.target;
            const formData = new FormData(form);
            
            fetch('/generate', {
                method: 'POST',
                body: formData
            })
            .then(response => response.json())
            .then(data => {
                const box = document.getElementById('password-box');
                const historyList = document.getElementById('history-list');
                box.textContent = '';
                
                let i = 0;
                const typeEffect = setInterval(() => {
                    if (i <= data.password.length) {
                        box.textContent = data.password.substring(0, i);
                        i++;
                    } else {
                        clearInterval(typeEffect);
                    }
                }, 50);

                const historyItem = document.createElement('div');
                historyItem.className = 'history-item';
                historyItem.innerHTML = '<span class="password-text" onclick="copyToClipboard(this)">'+data.password+'</span><span class="delete-btn" onclick="deletePassword(this)">&#215;</span>';
                historyList.insertBefore(historyItem, historyList.firstChild);
            });
        };
    </script>
</body>
</html>`

type PageData struct {
	Password string
	History  []string
}

var passwordHistory []string

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/generate", generateHandler)
	http.HandleFunc("/delete", deleteHandler)
	fmt.Println("Server started at http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("index").Parse(htmlTemplate))
	data := PageData{
		History: passwordHistory,
	}
	w.Header().Set("Content-Type", "text/html")
	tmpl.Execute(w, data)
}

func generateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	length, _ := strconv.Atoi(r.FormValue("length"))
	if length < 4 || length > 128 {
		length = 12
	}

	uppercase := r.FormValue("uppercase") == "on"
	lowercase := r.FormValue("lowercase") == "on"
	numbers := r.FormValue("numbers") == "on"
	special := r.FormValue("special") == "on"

	password := generatePassword(length, uppercase, lowercase, numbers, special)
	
	// Add to history
	passwordHistory = append([]string{password}, passwordHistory...)

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"password":"%s"}`, password)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Find and remove the password from history
	for i, pass := range passwordHistory {
		if pass == req.Password {
			passwordHistory = append(passwordHistory[:i], passwordHistory[i+1:]...)
			break
		}
	}

	w.WriteHeader(http.StatusOK)
}

func generatePassword(length int, uppercase, lowercase, numbers, special bool) string {
	var chars string
	if uppercase {
		chars += "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	}
	if lowercase {
		chars += "abcdefghijklmnopqrstuvwxyz"
	}
	if numbers {
		chars += "0123456789"
	}
	if special {
		chars += "!@#$%^&*()_+-=[]{}|;:,.<>?"
	}

	if chars == "" {
		chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	}

	password := make([]byte, length)
	for i := 0; i < length; i++ {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		password[i] = chars[n.Int64()]
	}

	return string(password)
}
