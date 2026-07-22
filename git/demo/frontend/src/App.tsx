import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import Header from './components/Header';
import Footer from './components/Footer';
import Home from './pages/Home';
import About from './pages/About';
import ProtectedRoute from './components/routes/ProtectedRoutes';
import AdminRoute from './components/routes/AdminRoute';
import DiscordCallback from './components/socials/DiscordCallback';
import N8NIntegrationSetup from './components/N8NIntegrationSetup';
import ClipAutomationPage from './components/ClipAutomationPage';
import SubscriptionDashboard from './components/SubscriptionDashboard';
import VoteSessionManager from './components/VoteSessionManager';
import WorkflowMarketplace from './components/WorkflowMarketplace';
import DashboardOverview from './components/userDashboard/DashboardOverview';
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
              {/* Public Routes */}
              <Route path='/' element={<Home />} />
              <Route path="/contact" element={<ContactPage />} />
              <Route path="/services" element={<ServicesPage />} />
              <Route path='/about' element={<About />} />

              {/* Legal */}
              <Route path="/legal/impressum" element={<Impressum />} />
              <Route path="/legal/datenschutz" element={<Datenschutz />} />
              <Route path="/legal/agb" element={<Agb />} />

              {/* Admin Routes - Protected for Admins only */}
              <Route
                path='/admin'
                element={
                  <AdminRoute>
                    <AdminDashboard />
                  </AdminRoute>
                }
              />

              {/* User Dashboard Routes - Protected for authenticated users */}
              <Route
                path='/dashboard'
                element={
                  <ProtectedRoute>
                    <Dashboard />
                  </ProtectedRoute>
                }
              >
                <Route index element={<DashboardOverview />} />
                <Route path='subscription' element={<SubscriptionDashboard />} />
                <Route path='n8n' element={<N8NIntegrationSetup />} />
                <Route path='clips' element={<ClipAutomationPage />} />
                <Route path='votes' element={<VoteSessionManager />} />
                <Route path='workflows' element={<WorkflowMarketplace />} />
              </Route>

              {/* Discord OAuth Callback */}
              <Route
                path='/discord/callback'
                element={
                  <ProtectedRoute>
                    <DiscordCallback />
                  </ProtectedRoute>
                }
              />

              {/* Catch all - redirect to home */}
              <Route path='*' element={<Navigate to="/" replace />} />
            </Routes>
          </div>
        </main>

        <Footer />
      </div>
    </Router>
  );
};

export default App;