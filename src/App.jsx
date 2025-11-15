import { Route, Routes } from "react-router-dom";
import TestRamzan from "./pages/test";
import Main from "./pages/main";
import BarcodeScannerPage from "./pages/barcode-scanner";
import AddNewMeal from "./pages/add-new-meal";
import BarcodeScannerPage2 from "./pages/barcode-scanner-2";

const App = () => (
  <Routes>
    <Route path="/" element={<Main />} />
    <Route
      path="/test"
      element={<AddNewMeal meal="Завтрак" date={"Четверг, 29 ноября"} />}
    />
    <Route path="/barcode" element={<BarcodeScannerPage2 />} />
  </Routes>
);

export default App;
