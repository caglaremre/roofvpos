const toast_success = document.getElementById('toast_success')
const toast_failure = document.getElementById('toast_failure')
const configButton = document.getElementById('configUpdate');
configButton.addEventListener('click', updateConfig)

async function getIPAddress() {
	try {
		const response = await fetch('https://ifconfig.me/ip');
		const data = await response.text();
		console.log("Your IP Address is:", data);
		document.getElementById('sale-cardholder-ip').value = data
		document.getElementById('sale-merchant-ip').value = data
		document.getElementById('sale-submerchant-ip').value = data
		document.getElementById('threeds-cardholder-ip').value = data
		document.getElementById('threeds-merchant-ip').value = data
		document.getElementById('threeds-submerchant-ip').value = data

	} catch (error) {
		console.error("Error fetching IP:", error);
	}
}

getIPAddress();

async function updateConfig() {
	const clientToken = document.getElementById('clientToken').value
	const secretKey = document.getElementById('secretKey').value
	const baseUrl = document.getElementById('baseUrl').value

	const response = await fetch('http://localhost:8080/config', {
		method: 'POST',
		headers: { 'content-type': 'application/json' },
		body: JSON.stringify({ clientToken: clientToken, secretKey: secretKey, baseUrl: baseUrl }),
	})
	if (response.ok) {
		const toastBootstrap = bootstrap.Toast.getOrCreateInstance(toast_success)
		toastBootstrap.show()
	} else {
		const toastBootstrap = bootstrap.Toast.getOrCreateInstance(toast_failure)
		const data = await response.json();
		document.getElementById('toast_failure_message').innerText = data.error
		toastBootstrap.show()
	}

}

document.getElementById('transaction-list').addEventListener('submit', async function (e) {
	if (e.target.name === 'order-void') {
		e.preventDefault();
		const form = e.target
		const orderID = form.querySelector('input[name="orderID"]').value;
		const processID = form.querySelector('input[name="processID"]').value;
		console.log("passing values to void canvas:", orderID, "Process:", processID);
		document.getElementById('void-order-id').value = orderID;
		document.getElementById('void-process-id').value = processID;
		document.getElementById('void-offcanvas-button').click()
	}
})

document.getElementById('transaction-list').addEventListener('submit', async function (e) {
	if (e.target.name === 'order-refund') {
		e.preventDefault();
		const form = e.target
		const orderID = form.querySelector('input[name="orderID"]').value;
		const processID = form.querySelector('input[name="processID"]').value;
		const amount = form.querySelector('input[name="amount"]').value;
		const pointAmount = form.querySelector('input[name="pointAmount"]').value;
		console.log("passing values to refund canvas:", orderID, "Process:", processID, "Amount", amount, "PointAmount", pointAmount);
		document.getElementById('refund-order-id').value = orderID;
		document.getElementById('refund-process-id').value = processID;
		document.getElementById('refund-amount').value = amount;
		document.getElementById('refund-point-amount').value = pointAmount;
		document.getElementById('refund-offcanvas-button').click()
	}
})

document.getElementById('transaction-list').addEventListener('submit', async function (e) {
	if (e.target.name === 'order-postsale') {
		e.preventDefault();
		const form = e.target
		const orderID = form.querySelector('input[name="orderID"]').value;
		const processID = form.querySelector('input[name="processID"]').value;
		const amount = form.querySelector('input[name="amount"]').value;
		console.log("passing values to postsale canvas:", orderID, "Process:", processID, "Amount", amount);
		document.getElementById('postsale-order-id').value = orderID;
		document.getElementById('postsale-process-id').value = processID;
		document.getElementById('postsale-amount').value = amount;
		document.getElementById('postsale-offcanvas-button').click()
	}
})

document.getElementById('transaction-list').addEventListener('click', async function (e) {
	const btn = e.target.closest('button')
	if (btn) {
		const offcanvasElement = document.getElementById('dynamicOffcanvas')
		const title = offcanvasElement.querySelector('#offcanvasTitle')
		const body = offcanvasElement.querySelector('#offcanvasContent')
		const orderid = btn.getAttribute('order-id')
		const response = await fetch(`http://localhost:8080/transaction/${orderid}`, {
			method: 'GET',
		})
		const html = await response.text()
		title.textContent = orderid
		body.innerHTML = html
	}
})