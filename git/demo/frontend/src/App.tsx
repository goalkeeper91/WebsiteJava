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
import SubathonPage from './components/userDashboard/SubathonPage';
import DiscordDashboard from './components/userDashboard/DiscordDashboard';
import TwitchLayout from './components/userDashboard/TwitchLayout';
import LiveDashboardHome from './components/userDashboard/LiveDashboardHome';
import CommandsDashboard from './components/userDashboard/CommandsDashboard';
import BuiltinCommandsInfo from './components/userDashboard/BuiltinCommandsInfo';
import VoteSessionManager from './components/VoteSessionManager';
import ScheduledMessagesDashboard from './components/userDashboard/ScheduledMessagesDashboard';
import AutomodDashboard from './components/userDashboard/AutomodDashboard';
import LoyaltyDashboard from './components/userDashboard/LoyaltyDashboard';
import GiveawaysDashboard from './components/userDashboard/GiveawaysDashboard';
import CS2CasterDashboard from './components/userDashboard/CS2CasterDashboard';
import TeamDashboard from './components/userDashboard/TeamDashboard';
import SettingsPage from './components/userDashboard/SettingsPage';
import { useTeam } from './context/TeamContext';
import AdminDashboard from './pages/admin/AdminDashboard';
import Impressum from "./pages/legal/Impressum";
import Datenschutz from "./pages/legal/Datenschutz";
import Agb from "./pages/legal/AGB";
import Widerruf from "./pages/legal/Widerruf";
import Cookies from "./pages/legal/Cookies";
import ContactPage from "./pages/ContactPage";
import ServicesPage from "./pages/ServicesPage";
import Pricing from "./pages/Pricing";
import Welcome from "./pages/Welcome";
import CancelContract from "./pages/CancelContract";
import Dashboard from "./pages/[user]/Dashboard";
import { LoginPopup } from './components/popup/LoginFailed';
import SubathonOverlay from './pages/SubathonOverlay';
import CS2MatchNotesPopup from './pages/CS2MatchNotesPopup';

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

// Team management only applies to your own channel, not one you're managing
// on someone else's behalf - deep-linking here while acting as another
// channel bounces back to the Twitch overview instead of rendering.
const TeamRoute = () => {
  const { actingAsChannel } = useTeam();
  return actingAsChannel ? <Navigate to="/dashboard/twitch" replace /> : <TeamDashboard />;
};

const App = () => {
  return (
    <Router>
      <Routes>
        {/* Bare routes - no header/footer chrome */}
        <Route path='/overlay/subathon' element={<SubathonOverlay />} />
        {/* Opened via window.open() as a real popup window from the CS2
            Caster Tools tab - deliberately bare so it reads as a compact
            moderation-card tool window, not a page inside the site. */}
        <Route path='/cs2/match-notes-popup' element={<CS2MatchNotesPopup />} />

        <Route element={<SiteLayout />}>
          {/* Public Routes */}
          <Route path='/' element={<Home />} />
          <Route path="/contact" element={<ContactPage />} />
          <Route path="/services" element={<ServicesPage />} />
          <Route path='/about' element={<About />} />
          <Route path="/pricing" element={<Pricing />} />
          <Route path="/welcome" element={<Welcome />} />
          {/* Kündigungsbutton nach § 312k BGB - permanently reachable from the
              footer on every page, not tucked away behind a dashboard login. */}
          <Route path="/vertrag-kuendigen" element={<CancelContract />} />

          {/* Legal */}
          <Route path="/legal/impressum" element={<Impressum />} />
          <Route path="/legal/datenschutz" element={<Datenschutz />} />
          <Route path="/legal/agb" element={<Agb />} />
          <Route path="/legal/widerruf" element={<Widerruf />} />
          <Route path="/legal/cookies" element={<Cookies />} />

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
            <Route index element={<Navigate to="twitch" replace />} />
            <Route path='settings' element={<SettingsPage />} />
            <Route path='subscription' element={<SubscriptionDashboard />} />
            <Route path='subathon' element={<SubathonPage />} />
            <Route path='n8n' element={<N8NIntegrationSetup />} />
            <Route path='clips' element={<ClipAutomationPage />} />
            <Route path='discord' element={<DiscordDashboard />} />

            <Route path='twitch' element={<TwitchLayout />}>
              <Route index element={<LiveDashboardHome />} />
              <Route path='commands' element={<CommandsDashboard />} />
              <Route path='builtin' element={<BuiltinCommandsInfo />} />
              <Route path='votes' element={<VoteSessionManager />} />
              <Route path='scheduled' element={<ScheduledMessagesDashboard />} />
              <Route path='automod' element={<AutomodDashboard />} />
              <Route path='loyalty' element={<LoyaltyDashboard />} />
              <Route path='giveaways' element={<GiveawaysDashboard />} />
              <Route path='cs2' element={<CS2CasterDashboard />} />
              <Route path='team' element={<TeamRoute />} />
            </Route>

            {/* Alte Bookmarks/Links auf die vorherige Tab-/Query-basierte Struktur */}
            <Route path='votes' element={<Navigate to="/dashboard/twitch/votes" replace />} />
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
