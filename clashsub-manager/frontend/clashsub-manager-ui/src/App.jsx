import './App.css'
import LoginPage from "./pages/login/Login.jsx";
import NotFoundPage from "./pages/error/NotFound.jsx";
import ServerErrorPage from "./pages/error/ServerError.jsx";
import MaintenancePage from "./pages/error/Maintenance.jsx";

function App() {

  return (
    <>
      <LoginPage></LoginPage>
      <NotFoundPage></NotFoundPage>
      <ServerErrorPage></ServerErrorPage>
      <MaintenancePage></MaintenancePage>

    </>)
}

export default App
