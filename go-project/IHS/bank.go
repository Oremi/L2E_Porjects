package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

const snapshotFile = "bank_snapshot.json"

type Customer struct {
	ID             int     `json:"id"`
	AccountNumber  string  `json:"account_number"`
	Name           string  `json:"name"`
	AccountBalance float64 `json:"account_balance"`
}

type Bank struct {
	Customers []Customer `json:"customers"`
	NextID    int        `json:"next_id"`
}

var bank Bank
var tmpl = template.Must(template.ParseFiles("templates/index.html"))

// ---------- Persistence ----------
func (b *Bank) Save() {
	data, _ := json.MarshalIndent(b, "", "  ")
	os.WriteFile(snapshotFile, data, 0600)
}

func Load() Bank {
	data, err := os.ReadFile(snapshotFile)
	if err != nil {
		return Bank{NextID: 1}
	}
	var b Bank
	if json.Unmarshal(data, &b) != nil {
		return Bank{NextID: 1}
	}
	fmt.Printf("Restored %d accounts\n", len(b.Customers))
	return b
}

// ---------- Helpers ----------
func genAccountNumber() string {
	return fmt.Sprintf("%04d-%04d", rand.Intn(10000), rand.Intn(10000))
}

func parseID(s string) (int, bool) {
	id, err := strconv.Atoi(s)
	return id, err == nil
}

func parseAmount(s string) (float64, bool) {
	a, err := strconv.ParseFloat(s, 64)
	return a, err == nil && a > 0
}

func redirect(w http.ResponseWriter, r *http.Request, errMsg string) {
	http.Redirect(w, r, "/?err="+url.QueryEscape(errMsg), http.StatusSeeOther)
}

func (b *Bank) AddCustomer(name string, deposit float64) Customer {
	c := Customer{
		ID:             b.NextID,
		AccountNumber:  genAccountNumber(),
		Name:           name,
		AccountBalance: deposit,
	}
	b.Customers = append(b.Customers, c)
	b.NextID++
	b.Save()
	return c
}

func (b *Bank) GetCustomer(id int) *Customer {
	for i := range b.Customers {
		if b.Customers[i].ID == id {
			return &b.Customers[i]
		}
	}
	return nil
}

func (b *Bank) Transfer(from, to int, amount float64) error {
	sender, receiver := b.GetCustomer(from), b.GetCustomer(to)
	if sender == nil || receiver == nil {
		return fmt.Errorf("invalid account")
	}
	if sender.ID == receiver.ID || sender.AccountBalance < amount {
		return fmt.Errorf("invalid transfer")
	}
	sender.AccountBalance -= amount
	receiver.AccountBalance += amount
	b.Save()
	return nil
}

type PageData struct {
	Customers []Customer
	Message   string
	Error     string
}

// ---------- Handlers ----------
func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	tmpl.Execute(w, PageData{
		Customers: bank.Customers,
		Message:   r.URL.Query().Get("msg"),
		Error:     r.URL.Query().Get("err"),
	})
}

func add(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	name := strings.TrimSpace(r.FormValue("name"))
	if name == "" {
		redirect(w, r, "Name cannot be blank")
		return
	}
	for _, ch := range name {
		if !((ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == ' ' || ch == '-') {
			redirect(w, r, "Name may contain letters, spaces, or dashes")
			return
		}
	}
	deposit, ok := parseAmount(r.FormValue("deposit"))
	if !ok || deposit < 0 {
		redirect(w, r, "Invalid deposit amount")
		return
	}
	newCust := bank.AddCustomer(name, deposit)
	msg := fmt.Sprintf("Account created for %s (ID: %d)", newCust.Name, newCust.ID)
	http.Redirect(w, r, "/?msg="+url.QueryEscape(msg), http.StatusSeeOther)
}

func transaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	id, ok := parseID(r.FormValue("id"))
	if !ok {
		redirect(w, r, "Invalid customer ID")
		return
	}
	amount, ok := parseAmount(r.FormValue("amount"))
	if !ok {
		redirect(w, r, "Amount must be positive")
		return
	}
	cust := bank.GetCustomer(id)
	if cust == nil {
		redirect(w, r, "Customer not found")
		return
	}
	action := r.FormValue("type")
	if action == "withdraw" && cust.AccountBalance < amount {
		redirect(w, r, "Insufficient funds")
		return
	}
	if action == "withdraw" {
		cust.AccountBalance -= amount
	} else {
		cust.AccountBalance += amount
	}
	bank.Save()
	msg := fmt.Sprintf("%s ₦%.2f to/from %s", strings.Title(action), amount, cust.Name)
	http.Redirect(w, r, "/?msg="+url.QueryEscape(msg), http.StatusSeeOther)
}

func transfer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	from, ok1 := parseID(r.FormValue("from_id"))
	to, ok2 := parseID(r.FormValue("to_id"))
	if !ok1 || !ok2 {
		redirect(w, r, "Invalid sender or receiver ID")
		return
	}
	amount, ok := parseAmount(r.FormValue("amount"))
	if !ok {
		redirect(w, r, "Amount must be positive")
		return
	}
	if err := bank.Transfer(from, to, amount); err != nil {
		redirect(w, r, err.Error())
		return
	}
	http.Redirect(w, r, "/?msg=Transfer+successful", http.StatusSeeOther)
}

// ---------- Main ----------
func main() {
	bank = Load()
	if len(bank.Customers) == 0 {
		bank.AddCustomer("Alice Smith", 1000)
		bank.AddCustomer("Bob Jones", 1500)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/add", add)
	mux.HandleFunc("/transaction", transaction)
	mux.HandleFunc("/transfer", transfer)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Server running on http://localhost:%s\n", port)
	http.ListenAndServe(":"+port, mux)
}
