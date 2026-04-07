config_button = document.getElementById('config-update')
config_button.addEventListener('click', updateConfig)

async function updateConfig() {
	const clientToken = document.getElementById('config-client-token').value
	const secretKey = document.getElementById('config-secret-key').value
	const baseUrl = document.getElementById('config-base-url').value

	const response = await fetch('http://localhost:8080/config', {
		method: 'POST',
		headers: { 'content-type': 'application/json' },
		body: JSON.stringify({ clientToken: clientToken, secretKey: secretKey, baseUrl: baseUrl }),
	})
	if (response.ok) {
		toast(true)
	} else {
		const data = await response.json();
		toast(false, data.error)
	}

}

function toast(state, message) {
	let alert
	if (state) {
		alert = document.querySelector('.alert-success')
	} else {
		alert = document.querySelector('.alert-error')
		document.getElementById('toast_failure_message').innerText = message
	}
	console.log(alert)
	alert.classList.remove('hidden')
	setTimeout(() => {
		alert.classList.add('hidden')
		config_button.classList.remove('btn-disabled')
	}, 3000);
}

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

document.getElementById('transaction-list').addEventListener('click', async function (e) {
	const btn = e.target.closest('button')
	if (btn) {
		const modal = document.getElementById('transaction-modal')
		const body = document.getElementById('transaction-model-content')
		const title = document.getElementById('transaction-modal-title')
		const orderid = btn.getAttribute('order-id')
		const response = await fetch(`http://localhost:8080/transaction/${orderid}`, {
			method: 'GET',
		})
		const html = await response.text()
		title.innerText = orderid
		body.innerHTML = html
		modal.showModal()
	}
})
