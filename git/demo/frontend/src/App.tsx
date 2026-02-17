import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import Header from './components/Header';
import Footer from './components/Footer';
import Home from './pages/Home';
import About from './pages/About';
import ProtectedRoute from './components/routes/ProtectedRoutes';
import DiscordCallback from './components/socials/DiscordCallback';
import AdminDashboard from './pages/admin/AdminDashboard';
import Impressum from "./pages/legal/Impressum";
import Datenschutz from "./pages/legal/Datenschutz";
import Agb from "./pages/legal/AGB";
import ContactPage from "./pages/ContactPage";
import ServicesPage from "./pages/ServicesPage";
import Dashboard from "./pages/[user]/Dashboard";
import { LoginPopup } from './components/popup/LoginFailed';

const App = () => {
  return (
    <Router>
      <div className='flex flex-col w-screen min-h-screen bg-black-100 text-white overflow-x-hidden'>
        <Header />

        <LoginPopup />

        <main className='flex-grow w-screen pt-17'>
          <div className='w-full max-w-screen'>
            <Routes>
              <Route path='/' element={<Home />} />
              <Route path="/contact" element={<ContactPage />} />
              <Route path="/services" element={<ServicesPage />} />
              <Route path='/about' element={<About />} />
              <Route path='/admin' element={
                    <ProtectedRoute>
                        <AdminDashboard />
                    </ProtectedRoute>
                  }
              />
              <Route path='/dashboard' element={
                  <ProtectedRoute>
                      <Dashboard />
                  </ProtectedRoute>
              } />
              <Route path='/discord/callback' element={
                  <ProtectedRoute>
                       <DiscordCallback />
                  </ProtectedRoute>
              } />
              <Route path="/legal/impressum" element={<Impressum />} />
              <Route path="/legal/datenschutz" element={<Datenschutz />} />
              <Route path="/legal/agb" element={<Agb />} />
            </Routes>
          </div>
        </main>

        <Footer />
      </div>
    </Router>
  );
};

export default App;
