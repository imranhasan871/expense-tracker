// Modal functions
function showBudgetModal() {
    document.getElementById('budgetModal').style.display = 'block';
    document.getElementById('budgetForm').reset();
    document.getElementById('formMessage').style.display = 'none';
}

function hideBudgetModal() {
    document.getElementById('budgetModal').style.display = 'none';
}

// Close modal when clicking outside
window.onclick = function(event) {
    const modal = document.getElementById('budgetModal');
    if (event.target === modal) {
        hideBudgetModal();
    }
}

// Save budget
async function saveBudget(event) {
    event.preventDefault();
    
    const formMessage = document.getElementById('formMessage');
    const form = event.target;
    const formData = new FormData(form);
    
    const data = {
        category_id: parseInt(formData.get('category_id')),
        year: parseInt(formData.get('year')),
        amount: parseFloat(formData.get('amount'))
    };
    
    try {
        const response = await fetch('/api/budgets', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(data)
        });
        
        const result = await response.json();
        
        if (response.ok) {
            formMessage.className = 'form-message success';
            formMessage.textContent = result.message || 'Budget saved successfully!';
            formMessage.style.display = 'block';
            
            setTimeout(() => {
                hideBudgetModal();
                fetchBudgets(); // Refresh the list and summary
            }, 1500);
        } else {
            formMessage.className = 'form-message error';
            formMessage.textContent = result.message || 'Failed to save budget';
            formMessage.style.display = 'block';
        }
    } catch (error) {
        formMessage.className = 'form-message error';
        formMessage.textContent = 'An error occurred. Please try again.';
        formMessage.style.display = 'block';
        console.error('Error:', error);
    }
}

// Fetch and display budgets
async function fetchBudgets() {
    const year = 2026; // Default for demo
    try {
        const response = await fetch(`/api/budgets?year=${year}`);
        const result = await response.json();
        
        if (response.ok && result.success) {
            updateDashboard(result.data.summary);
            renderBudgetsTable(result.data.budgets);
        }
    } catch (error) {
        console.error('Error fetching budgets:', error);
    }
}

function updateDashboard(summary) {
    if (!summary) return;
    
    // Update summary cards
    document.getElementById('totalAnnualBudget').textContent = formatCurrency(summary.total_annual_budget);
    document.getElementById('highestAllocation').textContent = formatCurrency(summary.highest_allocation);
    document.getElementById('savingsTarget').textContent = formatCurrency(summary.savings_target);
    document.getElementById('remainingBudget').textContent = formatCurrency(summary.remaining_budget);
}

function renderBudgetsTable(budgets) {
    const tableBody = document.getElementById('budgetsTableBody');
    if (!tableBody) return;

    if (!budgets || budgets.length === 0) {
        tableBody.innerHTML = `
            <tr>
                <td colspan="4" class="empty-state">No budget allocations found for this year.</td>
            </tr>
        `;
        return;
    }

    tableBody.innerHTML = budgets.map(b => `
        <tr>
            <td>
                <div class="category-info">
                    <span class="category-name font-bold">${escapeHtml(b.category_name)}</span>
                </div>
            </td>
            <td class="font-bold">${b.year}</td>
            <td>
                <div class="budget-allocation-cell">
                    <span class="amount-value font-bold">$${b.amount.toLocaleString(undefined, {minimumFractionDigits: 2})}</span>
                    <div class="progress-container" style="margin-top: 5px; max-width: 200px;">
                        <div class="progress-bar-bg">
                            <div class="progress-bar-fill" style="width: 0%"></div>
                        </div>
                        <div class="progress-header" style="justify-content: flex-start; gap: 10px; margin-top: 2px;">
                            <span class="spend-amount" style="font-size: 0.75rem;">Spend: $0.00</span>
                            <span class="percent" style="font-size: 0.75rem;">(0%)</span>
                        </div>
                    </div>
                </div>
            </td>
            <td class="text-right">
                <button class="btn-icon" onclick="editBudget(${b.category_id}, ${b.amount})" title="Edit Budget">
                    <svg aria-hidden="true" xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 256 256"><path fill="currentColor" d="M227.31 73.37L182.63 28.7a16 16 0 0 0-22.63 0L36.69 152A15.86 15.86 0 0 0 32 163.31V208a16 16 0 0 0 16 16h44.69a15.86 15.86 0 0 0 11.31-4.69L227.31 96a16 16 0 0 0 0-22.63M92.69 208H48v-44.69l88-88L180.69 120ZM192 108.69L147.31 64l24-24L216 84.69Z"/></svg>
                </button>
            </td>
        </tr>
    `).join('');
}


function editBudget(categoryId, amount) {
    showBudgetModal();
    document.getElementById('budgetCategory').value = categoryId;
    document.getElementById('budgetAmount').value = amount;
}

// Initial fetch
document.addEventListener('DOMContentLoaded', fetchBudgets);

// Utilities
function formatCurrency(value) {
    return '$' + value.toLocaleString(undefined, {minimumFractionDigits: 2, maximumFractionDigits: 2});
}

function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}
