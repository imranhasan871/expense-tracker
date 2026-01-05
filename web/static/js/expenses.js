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
    
    if (filters.startDate) params.append('start_date', filters.startDate);
    if (filters.endDate) params.append('end_date', filters.endDate);
    if (filters.categoryId) params.append('category_id', filters.categoryId);
    
    if (params.toString()) {
        url += '?' + params.toString();
    }
    
    try {
        const response = await fetch(url);
        const result = await response.json();
        
        if (response.ok && result.success) {
            renderExpensesTable(result.data);
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
        startDate: document.getElementById('filterStartDate').value,
        endDate: document.getElementById('filterEndDate').value,
        categoryId: document.getElementById('filterCategory').value
    };
    
    fetchExpenses(filters);
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
