package main

import (
	"fmt"
	"html/template"
	"math"
	"net/http"
	"strconv"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("index.html"))
		tmpl.Execute(w, nil)
	})

	http.HandleFunc("/calculate", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			return
		}

		fee, _ := strconv.ParseFloat(r.FormValue("fee"), 64)
		tip, _ := strconv.ParseFloat(r.FormValue("tip"), 64)
		quest, _ := strconv.ParseFloat(r.FormValue("quest"), 64)
		wallet, _ := strconv.ParseFloat(r.FormValue("wallet"), 64)

		income := ((fee + quest) * 0.98) + tip
		rawRemittance := wallet - income
		finalRemittance := math.Max(0, math.Ceil(rawRemittance))
		remaining := finalRemittance - rawRemittance

		fmt.Fprintf(w, `
			<div id="result" class="mt-6 animate-fade-in">
				<div class="bg-gray-50 rounded-xl p-4 border border-gray-100">
					<div class="flex justify-between mb-3">
						<span class="text-gray-600 font-medium">Your Net Income</span>
						<span class="text-gray-900 font-bold">₱%.2f</span>
					</div>
					<div class="flex justify-between items-center py-4 border-t border-gray-200">
						<span class="text-gray-900 font-bold text-lg">Total Remittance</span>
						<span class="text-[#D70F64] font-extrabold text-2xl">₱%.0f</span>
					</div>
					<div class="bg-blue-50 p-3 rounded-lg flex justify-between items-center">
						<span class="text-blue-700 text-sm font-medium">Wallet Credit After</span>
						<span class="text-blue-700 font-bold text-sm">₱%.2f</span>
					</div>
				</div>
				<p class="text-[10px] text-gray-400 mt-4 text-center italic">
					*Remittance is rounded up to the nearest Peso.
				</p>
			</div>
		`, income, finalRemittance, remaining)
	})

	fmt.Println("Server: http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
