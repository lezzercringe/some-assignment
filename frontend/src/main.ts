import { FetchOrder, type Order, type FetchError } from './fetch';
import './style.css';

document.querySelector<HTMLDivElement>('#app')!.innerHTML = `
 <div class="min-h-screen bg-gray-100 flex items-center justify-center p-4">
  <div class="w-full max-w-xl space-y-6">
    <div class="bg-white p-6 rounded-2xl shadow-md">
      <h2 class="text-xl font-semibold mb-4">Lookup Order</h2>
      <div id="validationError" class="text-red-600 text-sm mb-1 hidden">
      </div>
      <div class="flex items-center gap-3">
        <input
          id="orderIdInput"
          type="text"
          placeholder="Enter Order ID"
          class="flex-1 px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
        <button
          id="fetchOrderBtn"
          class="px-4 py-2 bg-blue-600 text-white font-semibold rounded-lg hover:bg-blue-700 disabled:bg-blue-200 transition"
        >
          Fetch
        </button>
      </div>
    </div>

    <div
      id="orderCard"
      class="bg-white p-6 rounded-2xl shadow-md hidden opacity-0 translate-y-2 transition-all duration-300"
    >
    </div>
  </div>
</div>
`;

const orderCard = document.querySelector('#orderCard')!;
const validationErrorField = document.querySelector<HTMLInputElement>('#validationError')!;
const idInput = document.querySelector<HTMLInputElement>('#orderIdInput')!;
const button = document.querySelector<HTMLButtonElement>('#fetchOrderBtn')!;

button.addEventListener("click", lookupHandler);
idInput.addEventListener('keydown', (event) => {
  if (event.key === 'Enter') {
    lookupHandler();
  }
});

let isLoading = false;

async function lookupHandler() {
  if (isLoading) { return; }

  const orderId = idInput.value.trim();
  if (orderId.length === 0) {
    setValidationError("order id cannot be empty");
    return;
  }

  setIsLoading(true)
  setValidationError(null);
  const result = await FetchOrder(orderId);
  setIsLoading(false)

  if (!result.ok) {
    showCard(renderFetchError(result.error));
    return;
  }

  showCard(renderOrderCard(result.value));
}

function setIsLoading(isLoading: boolean) {
  isLoading = true;
  button.disabled = isLoading;
  window.document.body.style["background"] = "red";
}

function setValidationError(message: string | null) {
  validationErrorField.textContent = message ?? '';
  validationErrorField.classList.toggle('hidden', !message);
}

function showCard(content :string) {
  orderCard.innerHTML = content;
  orderCard.classList.remove('hidden');

  requestAnimationFrame(() => {
    orderCard.classList.remove('opacity-0', 'translate-y-2');
  });
}

function renderFetchError(error: FetchError): string {
  if (error.type == "api") {
    return error.message;
  }

  return `Unexpected ${error.type} error happened`
}


function renderOrderCard(order: Order): string {
  const createdDate = new Date(order.date_created).toLocaleString();
  const paymentDate = new Date(order.payment.payment_dt * 1000).toLocaleString();
  const itemsHtml = order.items.map(item => `
    <tr class="border-b last:border-none">
      <td class="p-2">${item.name}</td>
      <td class="p-2 text-right">${item.price}</td>
      <td class="p-2 text-right">${item.sale}%</td>
      <td class="p-2 text-right">${item.total_price}</td>
    </tr>
  `).join('');

  return `
    <h3 class="text-lg font-bold mb-4">Order Details</h3>
    <div class="space-y-4 text-sm text-gray-800 max-h-[400px] overflow-auto">
      <section>
        <h4 class="font-semibold mb-2">General</h4>
        <p><strong>Order UID:</strong> ${order.order_uid}</p>
        <p><strong>Track Number:</strong> ${order.track_number}</p>
        <p><strong>Entry:</strong> ${order.entry}</p>
        <p><strong>Date Created:</strong> ${createdDate}</p>
        <p><strong>Locale:</strong> ${order.locale}</p>
        <p><strong>Delivery Service:</strong> ${order.delivery_service}</p>
      </section>

      <section>
        <h4 class="font-semibold mb-2">Delivery Info</h4>
        <p><strong>Name:</strong> ${order.delivery.name}</p>
        <p><strong>Phone:</strong> ${order.delivery.phone}</p>
        <p><strong>Address:</strong> ${order.delivery.address}, ${order.delivery.city}, ${order.delivery.region}, ${order.delivery.zip}</p>
        <p><strong>Email:</strong> ${order.delivery.email}</p>
      </section>

      <section>
        <h4 class="font-semibold mb-2">Payment</h4>
        <p><strong>Transaction ID:</strong> ${order.payment.transaction}</p>
        <p><strong>Request ID:</strong> ${order.payment.request_id}</p>
        <p><strong>Provider:</strong> ${order.payment.provider}</p>
        <p><strong>Amount:</strong> $${order.payment.amount}</p>
        <p><strong>Currency:</strong> ${order.payment.currency}</p>
        <p><strong>Payment Date:</strong> ${paymentDate}</p>
        <p><strong>Bank:</strong> ${order.payment.bank}</p>
        <p><strong>Delivery Cost:</strong> $${order.payment.delivery_cost}</p>
        <p><strong>Goods Total:</strong> $${order.payment.goods_total}</p>
        <p><strong>Custom Fee:</strong> $${order.payment.custom_fee}</p>
      </section>

      <section>
        <h4 class="font-semibold mb-2">Items</h4>
        <table class="w-full border-collapse text-left text-sm">
          <thead>
            <tr class="border-b bg-gray-50">
              <th class="p-2">Name</th>
              <th class="p-2 text-right">Price</th>
              <th class="p-2 text-right">Sale</th>
              <th class="p-2 text-right">Total Price</th>
            </tr>
          </thead>
          <tbody>
            ${itemsHtml}
          </tbody>
        </table>
      </section>
    </div>
  `;
}
