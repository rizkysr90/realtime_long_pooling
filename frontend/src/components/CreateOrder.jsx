import { useState } from "react";

// CreateOrder Component
const CreateOrder = () => {
  const [restaurantId, setRestaurantId] = useState(1);
  const [error, setError] = useState(null);
  const [success, setSuccess] = useState(false);

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
      setSuccess(true);
      setError(null);

      // Reset success message after 3 seconds
      setTimeout(() => setSuccess(false), 3000);
    } catch (err) {
      setError(err.message);
      setSuccess(false);
    }
  };

  return (
    <div className="p-4 max-w-2xl mx-auto">
      <h1 className="text-2xl font-bold mb-6">Create New Order</h1>

      <div className="space-y-4">
        <div className="space-y-2">
          <label className="block text-sm font-medium">Restaurant ID</label>
          <input
            type="number"
            value={restaurantId}
            onChange={(e) => setRestaurantId(Number(e.target.value))}
            className="border rounded px-3 py-2 w-full"
            min="1"
          />
        </div>

        <button
          onClick={createOrder}
          className="w-full bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600"
        >
          Create Order
        </button>

        {error && (
          <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded">
            {error}
          </div>
        )}

        {success && (
          <div className="bg-green-100 border border-green-400 text-green-700 px-4 py-3 rounded">
            Order created successfully!
          </div>
        )}
      </div>
    </div>
  );
};

export default CreateOrder;
