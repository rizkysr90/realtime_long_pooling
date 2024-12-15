import { useEffect, useState } from "react";

const OrdersComponent = () => {
  const [orders, setOrders] = useState([]);
  const [restaurantId, setRestaurantId] = useState(1);
  const [isPolling, setIsPolling] = useState(false);
  const [error, setError] = useState(null);

  useEffect(() => {
    let timeoutId;

    const pollOrders = async () => {
      try {
        console.log("Polling for orders...");
        const response = await fetch(
          `http://localhost:8080/poll/orders/${restaurantId}`
        );

        if (!response.ok) throw new Error("Polling failed");

        const data = await response.json();
        console.log("Received data:", data);

        if (data.status !== "timeout") {
          setOrders((prev) => [...prev, data]);
        }

        // Continue polling if still enabled
        if (isPolling) {
          timeoutId = setTimeout(pollOrders, 1000); // Small delay between polls
        }
      } catch (err) {
        console.error("Polling error:", err);
        setError(err.message);
        setIsPolling(false);
      }
    };

    if (isPolling) {
      pollOrders();
    }

    // Cleanup
    return () => {
      if (timeoutId) clearTimeout(timeoutId);
    };
  }, [isPolling, restaurantId]);

  const createOrder = async () => {
    try {
      const response = await fetch("http://localhost:8080/orders", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          restaurant_id: restaurantId,
          status: "new",
        }),
      });

      if (!response.ok) throw new Error("Failed to create order");

      const newOrder = await response.json();
      setOrders((prev) => [...prev, newOrder]);
    } catch (err) {
      setError(err.message);
    }
  };

  return (
    <div className="p-4 max-w-2xl mx-auto">
      <div className="mb-6 space-y-4">
        <h1 className="text-2xl font-bold">Restaurant Orders Dashboard</h1>

        <div className="flex gap-4 items-center">
          <input
            type="number"
            value={restaurantId}
            onChange={(e) => setRestaurantId(Number(e.target.value))}
            className="border rounded px-3 py-2 w-24"
            min="1"
          />

          <button
            onClick={createOrder}
            className="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600"
          >
            Create Order
          </button>

          <button
            onClick={() => setIsPolling(!isPolling)}
            className={`px-4 py-2 rounded ${
              isPolling
                ? "bg-red-500 hover:bg-red-600"
                : "bg-green-500 hover:bg-green-600"
            } text-white`}
          >
            {isPolling ? "Stop Polling" : "Start Polling"}
          </button>
        </div>

        {error && (
          <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded">
            {error}
          </div>
        )}
      </div>

      <div className="space-y-4">
        <h2 className="text-xl font-semibold">Orders</h2>
        {orders.length === 0 ? (
          <p className="text-gray-500">No orders yet</p>
        ) : (
          <div className="space-y-2">
            {orders.map((order) => (
              <div
                key={order.id}
                className="border rounded p-4 shadow-sm hover:shadow-md transition-shadow"
              >
                <div className="flex justify-between">
                  <span className="font-medium">Order #{order.id}</span>
                  <span
                    className={`px-2 py-1 rounded-full text-sm ${
                      order.status === "new"
                        ? "bg-green-100 text-green-800"
                        : "bg-gray-100"
                    }`}
                  >
                    {order.status}
                  </span>
                </div>
                <div className="text-sm text-gray-600 mt-2">
                  Created: {new Date(order.created_at).toLocaleString()}
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
};

export default OrdersComponent;
