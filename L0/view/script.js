const searchInput = document.getElementById("searchInput");

function PostOrders() {
  fetch(`/order/fake`, {
    method: "POST"
  })
    .then(response => response.json())
    .catch(error => {
      console.error('Error:', error)
    });
}

function GetOrderById() {
  var input = searchInput.value

  if (input.trim() == "") {
    alert("No order id");
    return;
  }

  var orderId = String(input)
  fetch(`/order/${orderId}`, {
    method: "GET"
  })
    .then(response => response.json())
    .then(response => {
      if (response.error) {
        console.error('Error: ', response.message);
      } else {
        const order = response;
        clearTables();
        fillOrderInfo(order);
      }
    })
    .catch(error => {
      console.error('Error: ', error);
    });

}

function fillOrderInfo(order) {
  var ignoredFields = ['delivery', 'payment', 'items'];
  fillObject('order', ignoredFields, order);
  fillObject('delivery', [], order.delivery);
  fillObject('payment', [], order.payment);

  const itemsBody = document.getElementById('tableItems').querySelector('tbody');
  order.items.forEach((prop) => {

    const row = document.createElement('tr');
    row.insertCell().textContent = prop.chrt_id;
    row.insertCell().textContent = prop.track_number;
    row.insertCell().textContent = prop.price;
    row.insertCell().textContent = prop.rid;
    row.insertCell().textContent = prop.name;
    row.insertCell().textContent = prop.sale;
    row.insertCell().textContent = prop.size;
    row.insertCell().textContent = prop.total_price;
    row.insertCell().textContent = prop.nm_id;
    row.insertCell().textContent = prop.brand;
    row.insertCell().textContent = prop.status;
    itemsBody.appendChild(row);
  });
}

function fillObject(tableNmae, ignoredFields, entry) {
  var table = document.getElementById(tableNmae).querySelector('tbody');
  row = document.createElement('tr');
  Object.entries(entry).forEach(([key, value]) => {

    if (!ignoredFields.includes(key)) {
      const cell = document.createElement('td');
      cell.textContent = value;
      row.appendChild(cell);
    }
  });
  table.appendChild(row);
}

function clearTables() {
  document.getElementById('order').querySelector('tbody').innerHTML = '';
  document.getElementById('delivery').querySelector('tbody').innerHTML = '';
  document.getElementById('payment').querySelector('tbody').innerHTML = '';
  document.getElementById('tableItems').querySelector('tbody').innerHTML = '';
}