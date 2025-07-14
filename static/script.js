async function getOrderInfo() {
    const orderUidInput = document.getElementById('orderUid');
    const resultDiv = document.getElementById('result');
    const orderUid = orderUidInput.value.trim();

    resultDiv.innerHTML = '';

    if (!orderUid) {
        resultDiv.innerHTML = '<p class="error">Please enter an Order UID.</p>';
        return;
    }

    const backendUrl = `http://localhost:8081/order/${orderUid}`;

    resultDiv.innerHTML = '<p>Loading...</p>';

    try {
        const response = await fetch(backendUrl);

        if (response.ok) {
            const data = await response.json();
            displayOrderDetails(data, resultDiv);
        } else if (response.status === 404) {
            resultDiv.innerHTML = '<p class="error">Order not found for UID: ' + orderUid + '</p>';
        } else {
            const errorText = await response.text();
            resultDiv.innerHTML = '<p class="error">Error fetching data: ' + response.status + ' ' + response.statusText + '<br>Details: ' + errorText + '</p>';
        }
    } catch (error) {
        console.error('Fetch error:', error);
        resultDiv.innerHTML = '<p class="error">Network error or API is unreachable. Check console for details.</p>';
    }
}

function displayOrderDetails(order, targetElement) {
    let html = '<h2>Order Details:</h2>';
    html += '<div class="order-details">';

    html += '<div class="order-section">';
    html += '<h4>General Information</h4>';
    html += `<div class="order-prop"><strong>Order UID:</strong> ${order.OrderUID || 'N/A'}</div>`;
    html += `<div class="order-prop"><strong>Track Number:</strong> ${order.TrackNumber || 'N/A'}</div>`;
    html += `<div class="order-prop"><strong>Entry:</strong> ${order.Entry || 'N/A'}</div>`;
    html += `<div class="order-prop"><strong>Delivery Service:</strong> ${order.DeliveryService || 'N/A'}</div>`;
    html += `<div class="order-prop"><strong>Amount:</strong> ${order.Payment ? (order.Payment.Amount || 'N/A') : 'N/A'} ${order.Payment ? (order.Payment.Currency || 'N/A') : 'N/A'}</div>`; // Изменено: используем Payment.Amount и Payment.Currency
    html += `<div class="order-prop"><strong>Payment Method:</strong> ${order.Payment ? (order.Payment.Provider || 'N/A') : 'N/A'}</div>`;
    html += `<div class="order-prop"><strong>Locale:</strong> ${order.Locale || 'N/A'}</div>`;
    html += `<div class="order-prop"><strong>Internal Signature:</strong> ${order.InternalSignature || 'N/A'}</div>`;
    html += `<div class="order-prop"><strong>Customer ID:</strong> ${order.CustomerID || 'N/A'}</div>`;
    html += `<div class="order-prop"><strong>Shard Key:</strong> ${order.Shardkey || 'N/A'}</div>`;
    html += `<div class="order-prop"><strong>Sm ID:</strong> ${order.SmID || 'N/A'}</div>`;
    html += `<div class="order-prop"><strong>Date Created:</strong> ${order.DateCreated || 'N/A'}</div>`;
    html += `<div class="order-prop"><strong>OOF Shard:</strong> ${order.OofShard || 'N/A'}</div>`;
    html += '</div>';

    if (order.Payment) {
        html += '<div class="order-section">';
        html += '<h4>Payment Information</h4>';
        html += `<div class="order-prop"><strong>Transaction:</strong> ${order.Payment.Transaction || 'N/A'}</div>`;
        html += `<div class="order-prop"><strong>Request ID:</strong> ${order.Payment.RequestID || 'N/A'}</div>`;
        html += `<div class="order-prop"><strong>Currency:</strong> ${order.Payment.Currency || 'N/A'}</div>`;
        html += `<div class="order-prop"><strong>Provider:</strong> ${order.Payment.Provider || 'N/A'}</div>`;
        html += `<div class="order-prop"><strong>Amount:</strong> ${order.Payment.Amount || 'N/A'}</div>`;
        html += `<div class="order-prop"><strong>Payment DT:</strong> ${order.Payment.PaymentDT || 'N/A'}</div>`;
        html += `<div class="order-prop"><strong>Bank:</strong> ${order.Payment.Bank || 'N/A'}</div>`;
        html += `<div class="order-prop"><strong>Delivery Cost:</strong> ${order.Payment.DeliveryCost || 'N/A'}</div>`;
        html += `<div class="order-prop"><strong>Goods Total:</strong> ${order.Payment.GoodsTotal || 'N/A'}</div>`;
        html += `<div class="order-prop"><strong>Custom Fee:</strong> ${order.Payment.CustomFee || 'N/A'}</div>`;
        html += '</div>';
    }

    if (order.DeliveryInfo) {
        html += '<div class="order-section">';
        html += '<h4>Delivery Information</h4>';
        html += `<div class="order-prop"><strong>Name:</strong> ${order.DeliveryInfo.Name || 'N/A'}</div>`;
        html += `<div class="order-prop"><strong>Phone:</strong> ${order.DeliveryInfo.Phone || 'N/A'}</div>`;
        html += `<div class="order-prop"><strong>Zip:</strong> ${order.DeliveryInfo.Zip || 'N/A'}</div>`;
        html += `<div class="order-prop"><strong>City:</strong> ${order.DeliveryInfo.City || 'N/A'}</div>`;
        html += `<div class="order-prop"><strong>Address:</strong> ${order.DeliveryInfo.Address || 'N/A'}</div>`;
        html += `<div class="order-prop"><strong>Region:</strong> ${order.DeliveryInfo.Region || 'N/A'}</div>`;
        html += `<div class="order-prop"><strong>Email:</strong> ${order.DeliveryInfo.Email || 'N/A'}</div>`;
        html += '</div>';
    }

    if (order.Items && order.Items.length > 0) {
        html += '<div class="order-section">';
        html += '<h4>Items</h4>';
        html += '<ul class="order-item-list">';
        order.Items.forEach(item => {
            html += `<li>
                        <strong>Name:</strong> ${item.Name || 'N/A'}<br>
                        <strong>Quantity:</strong> ${item.ChrtID || 'N/A'}<br>  <strong>Price:</strong> ${item.Price || 'N/A'} (Sale: ${item.Sale || 'N/A'}%, Size: ${item.Size || 'N/A'})<br>
                        <strong>Total:</strong> ${item.TotalPrice || 'N/A'}<br>
                        <strong>Brand:</strong> ${item.Brand || 'N/A'}<br>
                        <strong>Status:</strong> ${item.Status || 'N/A'}
                      </li>`;
        });
        html += '</ul>';
        html += '</div>';
    } else {
        html += '<div class="order-section"><p>No items found for this order.</p></div>';
    }

    html += '</div>';
    targetElement.innerHTML = html;
}