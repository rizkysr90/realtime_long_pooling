import { BrowserRouter, Route, Routes } from "react-router-dom";
import CreateOrder from "./components/CreateOrder";
import RestaurantDashboard from "./components/RestaurantDashboard";

function App() {
  return (
    <BrowserRouter>
      <div>
        <Routes>
          <Route path="/" element={<CreateOrder />} />
          <Route path="/dashboard" element={<RestaurantDashboard />} />
        </Routes>
      </div>
    </BrowserRouter>
  );
}
export default App;
