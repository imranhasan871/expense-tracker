// Modal functions
function showExpenseModal() {
    const modal = document.getElementById('expenseModal');
    if (!modal) return;
    modal.style.display = 'block';
    document.getElementById('expenseForm').reset();
    document.getElementById('expenseDate').valueAsDate = new Date(); // Reset to today
    document.getElementById('formMessage').style.display = 'none';
}

function hideExpenseModal() {
    const modal = document.getElementById('expenseModal');
    if (modal) modal.style.display = 'none';
}

// Close modal when clicking outside
window.onclick = function(event) {
    const modal = document.getElementById('expenseModal');
    if (event.target === modal) {
        hideExpenseModal();
    }
}

// Save expense
async function saveExpense(event) {
    event.preventDefault();
    
    const formMessage = document.getElementById('formMessage');
    const form = event.target;
    const formData = new FormData(form);
    
    const data = {
        amount: parseFloat(formData.get('amount')),
        expense_date: formData.get('date'),
        category_id: parseInt(formData.get('category_id')),
        remarks: formData.get('remarks')
    };
    
    try {
        const response = await fetch('/api/expenses', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(data)
        });
        
        const result = await response.json();
        
        if (response.ok) {
            formMessage.className = 'form-message success';
            formMessage.textContent = result.message || 'Expense recorded successfully!';
            formMessage.style.display = 'block';
            
            setTimeout(() => {
                hideExpenseModal();
                fetchExpenses(); // Refresh the list
            }, 1500);
        } else {
            formMessage.className = 'form-message error';
            formMessage.textContent = result.message || 'Failed to record expense';
            formMessage.style.display = 'block';
        }
    } catch (error) {
        formMessage.className = 'form-message error';
        formMessage.textContent = 'An error occurred. Please try again.';
        formMessage.style.display = 'block';
        console.error('Error:', error);
    }
}

// Fetch and display expenses
async function fetchExpenses(filters = {}) {
    let url = '/api/expenses';
    const params = new URLSearchParams();
    
    if (filters.search) params.append('search', filters.search);
    if (filters.startDate) params.append('start_date', filters.startDate);
    if (filters.endDate) params.append('end_date', filters.endDate);
    if (filters.categoryId) params.append('category_id', filters.categoryId);
    if (filters.minAmount) params.append('min_amount', filters.minAmount);
    if (filters.maxAmount) params.append('max_amount', filters.maxAmount);
    
    if (params.toString()) {
        url += '?' + params.toString();
    }
    
    try {
        const response = await fetch(url);
        const result = await response.json();
        
        if (response.ok && result.success) {
            renderExpensesTable(result.data);
        } else if (!response.ok) {
            alert(result.message || 'Failed to fetch expenses');
        }
    } catch (error) {
        console.error('Error fetching expenses:', error);
    }
}

function renderExpensesTable(expenses) {
    const tableBody = document.getElementById('expensesTableBody');
    if (!tableBody) return;

    if (!expenses || expenses.length === 0) {
        tableBody.innerHTML = `
            <tr>
                <td colspan="5" class="empty-state">No expense records found matching your filters.</td>
            </tr>
        `;
        return;
    }

    tableBody.innerHTML = expenses.map(e => {
        const date = new Date(e.expense_date).toLocaleDateString('en-US', {
            year: 'numeric',
            month: 'short',
            day: '2-digit'
        });

        return `
            <tr>
                <td>${date}</td>
                <td class="font-bold">${escapeHtml(e.remarks)}</td>
                <td><span class="category-tag">${escapeHtml(e.category_name)}</span></td>
                <td><span class="expense-amount">$${e.amount.toLocaleString(undefined, {minimumFractionDigits: 2})}</span></td>
                <td class="text-right">
                    <div class="action-group">
                        <button class="btn-icon" onclick="deleteExpense(${e.id})" aria-label="Delete Expense" title="Delete Expense">
                            <svg aria-hidden="true" xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 256 256"><path fill="currentColor" d="M216 48h-40v-8a24 24 0 0 0-24-24h-48a24 24 0 0 0-24 24v8H40a8 8 0 0 0 0 16h8v144a16 16 0 0 0 16 16h128a16 16 0 0 0 16-16V64h8a8 8 0 0 0 0-16M96 40a8 8 0 0 1 8-8h48a8 8 0 0 1 8 8v8H96Zm96 168H64V64h128Zm-80-104v64a8 8 0 0 1-16 0v-64a8 8 0 0 1 16 0m48 0v64a8 8 0 0 1-16 0v-64a8 8 0 0 1 16 0"/></svg>
                        </button>
                    </div>
                </td>
            </tr>
        `;
    }).join('');
}

// Apply filters
function applyFilters() {
    const filters = {
        search: document.getElementById('filterSearch').value.trim(),
        startDate: document.getElementById('filterStartDate').value,
        endDate: document.getElementById('filterEndDate').value,
        categoryId: document.getElementById('filterCategory').value,
        minAmount: document.getElementById('filterMinAmount').value,
        maxAmount: document.getElementById('filterMaxAmount').value
    };
    
    // Client-side validation
    if (filters.startDate && filters.endDate && filters.startDate > filters.endDate) {
        alert('Start date must be before or equal to end date');
        return;
    }
    if (filters.minAmount && filters.maxAmount && parseFloat(filters.minAmount) > parseFloat(filters.maxAmount)) {
        alert('Minimum amount must be less than or equal to maximum amount');
        return;
    }
    
    fetchExpenses(filters);
}

// Clear all filters
function clearFilters() {
    document.getElementById('filterSearch').value = '';
    document.getElementById('filterStartDate').value = '';
    document.getElementById('filterEndDate').value = '';
    document.getElementById('filterCategory').value = '';
    document.getElementById('filterMinAmount').value = '';
    document.getElementById('filterMaxAmount').value = '';
    fetchExpenses();
}

// Delete expense
async function deleteExpense(id) {
    if (!confirm('Are you sure you want to delete this expense record?')) {
        return;
    }

    try {
        const response = await fetch(`/api/expenses/${id}`, {
            method: 'DELETE'
        });
        
        const result = await response.json();
        
        if (response.ok && result.success) {
            fetchExpenses(); // Refresh list
        } else {
            alert(result.message || 'Failed to delete expense');
        }
    } catch (error) {
        console.error('Error deleting expense:', error);
        alert('An error occurred. Please try again.');
    }
}

// Initial fetch
document.addEventListener('DOMContentLoaded', () => {
    fetchExpenses();
});

// Utilities
function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

// Budget Checking Logic
document.addEventListener('DOMContentLoaded', () => {
    const dateInput = document.getElementById('expenseDate');
    const categoryInput = document.getElementById('expenseCategory');
    const amountInput = document.getElementById('expenseAmount');

    if (dateInput && categoryInput && amountInput) {
        // Debounce the call slightly for input events
        const inputs = [dateInput, categoryInput, amountInput];
        inputs.forEach(elem => {
            elem.addEventListener('change', checkBudgetStatus);
            if (elem === amountInput) {
                elem.addEventListener('input', checkBudgetStatus);
            }
        });
    }
});

let budgetFetchTimeout;

async function checkBudgetStatus() {
    const dateVal = document.getElementById('expenseDate').value;
    const categoryId = document.getElementById('expenseCategory').value;
    const amountVal = parseFloat(document.getElementById('expenseAmount').value) || 0;
    const statusDiv = document.getElementById('budgetStatus');
    const uploadDiv = document.getElementById('approvalUpload');
    const fileInput = document.getElementById('approvedScan');

    if (!dateVal || !categoryId) {
        if (statusDiv) statusDiv.style.display = 'none';
        if (uploadDiv) uploadDiv.style.display = 'none';
        if (fileInput) fileInput.required = false;
        return;
    }

    // Clear previous timeout to debounce
    if (budgetFetchTimeout) clearTimeout(budgetFetchTimeout);

    budgetFetchTimeout = setTimeout(async () => {
        const year = new Date(dateVal).getFullYear();

        try {
            const response = await fetch(`/api/budgets/status?category_id=${categoryId}&year=${year}`);
            const result = await response.json();
            
            if (response.ok && result.success) {
                const data = result.data;
                const allocated = data.allocated;
                
                if (allocated === 0) {
                    // No budget info
                     statusDiv.style.display = 'none';
                     return;
                }

                const spentSoFar = data.spent;
                const currentTotal = spentSoFar + amountVal;
                
                let percent = 0;
                if (allocated > 0) {
                    percent = (currentTotal / allocated) * 100;
                }

                // Display Status
                statusDiv.style.display = 'block';
                const remaining = allocated - currentTotal;
                
                const barColor = percent >= 90 ? '#ef4444' : '#10b981'; // Red or Green
                const textColor = percent >= 90 ? '#dc2626' : '#047857';

                statusDiv.innerHTML = `
                    <div style="display: flex; justify-content: space-between; margin-bottom: 6px; font-weight: 500; color: #334155;">
                        <span style="color: ${textColor}">Usage: ${percent.toFixed(1)}%</span>
                        <span>Remaining: $${remaining.toLocaleString(undefined, {minimumFractionDigits: 2, maximumFractionDigits: 2})}</span>
                    </div>
                    <div style="width: 100%; height: 8px; background-color: #e2e8f0; border-radius: 4px; overflow: hidden;">
                        <div style="width: ${Math.min(percent, 100)}%; height: 100%; background-color: ${barColor}; transition: width 0.3s ease;"></div>
                    </div>
                    <div style="text-align: right; font-size: 0.75rem; color: #64748b; margin-top: 4px;">
                        Allocated: $${allocated.toLocaleString()}
                    </div>
                `;

                // Check 90% threshold for upload
                if (percent >= 90) {
                    if (uploadDiv.style.display === 'none') {
                        uploadDiv.style.display = 'block';
                        // Add a small animation effect
                        uploadDiv.animate([
                            { opacity: 0, transform: 'translateY(-10px)' },
                            { opacity: 1, transform: 'translateY(0)' }
                        ], { duration: 300 });
                    }
                    fileInput.required = true;
                } else {
                    uploadDiv.style.display = 'none';
                    fileInput.required = false;
                }
            }
        } catch (e) {
            console.error("Failed to fetch budget status", e);
        }
    }, 100); // 100ms debounce
}
