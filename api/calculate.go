package handler

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	// Only allow POST requests from HTMX
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse inputs from the form
	fee, _ := strconv.ParseFloat(r.FormValue("fee"), 64)
	tip, _ := strconv.ParseFloat(r.FormValue("tip"), 64)
	quest, _ := strconv.ParseFloat(r.FormValue("quest"), 64)
	wallet, _ := strconv.ParseFloat(r.FormValue("wallet"), 64)

	// Calculation Logic
	// Income = (Fee + Quest - 2% tax) + Tip
	income := ((fee + quest) * 0.98) + tip
	
	// Remittance = Wallet Balance - Net Income
	rawRemittance := wallet - income
	
	// Foodpanda remittance doesn't accept centavos. 
	// We round UP to ensure your app balance stays at 0.00 or positive.
	finalRemittance := math.Max(0, math.Ceil(rawRemittance))
	
	// Remaining in app = What you paid - What you owed
	remaining := finalRemittance - rawRemittance

	// Return the result fragment
	fmt.Fprintf(w, `
		<div id="result" class="mt-6 animate-fade-in space-y-3">
			<div class="bg-[#F7F7F7] rounded-xl p-4 border border-gray-100">
				<div class="flex justify-between items-center mb-1">
					<span class="text-gray-500 text-xs font-semibold uppercase tracking-tight">Your Net Income</span>
					<span class="text-gray-900 font-bold text-lg">₱%.2f</span>
				</div>
				
				<div class="flex justify-between items-center py-4 border-t border-gray-200 mt-2">
					<span class="text-gray-900 font-bold text-base">TO REMIT</span>
					<span class="text-[#D70F64] font-black text-3xl">₱%.0f</span>
				</div>

				<div class="bg-blue-50 p-3 rounded-lg flex justify-between items-center">
					<span class="text-blue-700 text-[10px] font-bold uppercase">App Wallet After</span>
					<span class="text-blue-700 font-bold">₱%.2f</span>
				</div>
			</div>
			<p class="text-[10px] text-gray-400 text-center italic">
				Remittance rounded up. The excess centavos will remain as credit in your rider app.
			</p>
		</div>
	`, income, finalRemittance, remaining)
}
