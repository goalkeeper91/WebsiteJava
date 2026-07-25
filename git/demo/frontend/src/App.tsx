import { BrowserRouter as Router, Routes, Route, Navigate, Outlet } from 'react-router-dom';
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
import SubathonOverlay from './pages/SubathonOverlay';

// Wraps every normal page in the shared header/footer chrome. Routes that
// need to render bare (no chrome at all - e.g. the OBS overlay, which must
// be a transparent Browser Source, not a page inside the site) live outside
// this layout entirely instead of just hiding Header/Footer with CSS.
const SiteLayout = () => (
  <div className='grid grid-rows-[auto_1fr_auto] w-full min-h-screen bg-black-100 text-white overflow-x-hidden'>
    <Header />
    <LoginPopup />
    {/* min-w-0: grid items default to min-width:auto, so without this any
        sufficiently wide descendant (a table, a long unbreakable string...)
        pushes this track wider than the viewport - the overflow-x-hidden
        above then hides the resulting scrollbar but silently CLIPS that
        content instead of letting it shrink/wrap, which is invisible in a
        scrollWidth-based check but very much visible (and broken) on a
        real phone. */}
    <main className='w-full min-w-0 pt-17'>
      <div className='w-full max-w-full min-w-0'>
        <Outlet />
      </div>
    </main>
    <Footer />
  </div>
);

const App = () => {
  return (
    <Router>
      <Routes>
        {/* Bare routes - no header/footer chrome */}
        <Route path='/overlay/subathon' element={<SubathonOverlay />} />

        <Route element={<SiteLayout />}>
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
            {/* Umfragen sind jetzt Unterbereich des Twitch-Chatbot-Tabs */}
            <Route path='votes' element={<Navigate to="/dashboard?tab=twitch" replace />} />
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
        </Route>
      </Routes>
    </Router>
  );
};

export default App;
