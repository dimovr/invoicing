// This would be included in your static/js/invoice.js file
document.addEventListener('htmx:load', function() {
    // Listen for the custom event when line items are added or removed
    document.body.addEventListener('htmx:afterSwap', function(event) {
        // Check if we need to update invoice summary
        if (event.detail.target.closest('#lineItemsTable tbody')) {
            updateInvoiceSummary();
        }
    });

    // Listen for our custom event triggered by the server
    document.body.addEventListener('updateInvoiceSummary', function(event) {
        updateInvoiceSummary();
    });
});

function updateInvoiceSummary() {
    // Get all line item rows
    const lineItems = document.querySelectorAll('.line-item-row');
    
    // Calculate totals
    let subtotal = 0;
    let taxAmount = 0;
    let total = 0;
    
    lineItems.forEach(item => {
        subtotal += parseFloat(item.dataset.subtotal || 0);
        taxAmount += parseFloat(item.dataset.taxAmount || 0);
        total += parseFloat(item.dataset.total || 0);
    });
    
    // Update the summary display
    document.querySelector('#invoiceSummary strong:nth-child(1)').textContent = `Subtotal: $${subtotal.toFixed(2)}`;
    document.querySelector('#invoiceSummary strong:nth-child(2)').textContent = `Tax Amount: $${taxAmount.toFixed(2)}`;
    document.querySelector('#invoiceSummary strong:nth-child(3)').textContent = `Total: $${total.toFixed(2)}`;
    
    // Update hidden fields for form submission
    document.querySelector('#hiddenSubtotal').value = subtotal.toFixed(2);
    document.querySelector('#hiddenTaxAmount').value = taxAmount.toFixed(2);
    document.querySelector('#hiddenTotal').value = total.toFixed(2);
    
    // Enable/disable the save button based on whether there are line items
    document.querySelector('#saveInvoiceBtn').disabled = lineItems.length === 0;
    
    // Create hidden inputs for line items to be submitted with the form
    const existingHiddenInputs = document.querySelectorAll('.line-item-hidden-input');
    existingHiddenInputs.forEach(input => input.remove());
    
    // Add hidden inputs for each line item
    const form = document.querySelector('#invoiceCreateForm');
    lineItems.forEach((item, index) => {
        const itemIdInput = document.createElement('input');
        itemIdInput.type = 'hidden';
        itemIdInput.name = `LineItems[ItemID]`;
        itemIdInput.value = item.dataset.id;
        itemIdInput.className = 'line-item-hidden-input';
        form.appendChild(itemIdInput);
        
        const countInput = document.createElement('input');
        countInput.type = 'hidden';
        countInput.name = `LineItems[Count]`;
        countInput.value = item.dataset.count;
        countInput.className = 'line-item-hidden-input';
        form.appendChild(countInput);
    });
}