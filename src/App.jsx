import { Route, Routes } from "react-router-dom";
import TestRamzan from "./pages/test";
import Main from "./pages/main";
import BarcodeScannerPage from "./pages/barcode-scanner";
import AddNewMeal from "./pages/add-new-meal";

const App = () => (
  <Routes>
    <Route path="/" element={<Main />} />
    <Route
      path="/test"
      element={<AddNewMeal meal="Завтрак" date={"Четверг, 29 ноября"} />}
    />
    <Route path="/barcode" element={<BarcodeScannerPage />} />
  </Routes>
);

export default App;
