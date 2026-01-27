function setupTableSearch(tableId, searchInputId, searchColumnIndex = 0) {
    const searchInput = document.getElementById(searchInputId);
    const table = document.getElementById(tableId);

    if (!searchInput || !table) return;

    searchInput.addEventListener('input', function(e) {
        const searchTerm = e.target.value.toLowerCase();
        const rows = table.querySelectorAll('tbody tr');

        rows.forEach(row => {
            const cell = row.cells[searchColumnIndex];
            const cellText = cell.textContent.toLowerCase();
            row.style.display = cellText.includes(searchTerm) ? '' : 'none';
        });
    });
}

// Auto-initialize for common patterns
document.addEventListener('DOMContentLoaded', function() {
    // Items search - search column index 1 (Proizvod)
    if (document.getElementById('productSearch') && document.getElementById('itemsTable')) {
        setupTableSearch('itemsTable', 'productSearch', 1);
    }

    // Suppliers search - search column index 0 (Ime)
    if (document.getElementById('supplierSearch') && document.getElementById('suppliersTable')) {
        setupTableSearch('suppliersTable', 'supplierSearch', 0);
    }
});
