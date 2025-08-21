import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import Header from './components/Header';
import Footer from './components/Footer';
import Home from './pages/Home';
import About from './pages/About';
import AllVideos from './pages/AllVideos';
import ProtectedRoute from './components/routes/ProtectedRoutes';
import AdminDashboard from './pages/admin/AdminDashboard';
import Impressum from "./pages/legal/Impressum";
import Datenschutz from "./pages/legal/Datenschutz";
import Agb from "./pages/legal/AGB";

const App = () => {
  return (
    <Router>
      <div className='flex flex-col w-screen h-screen bg-black-100 text-white'>
        <Header />

        <main className='flex-grow pt-19 w-full px-4 py-6 sm:px-6 lg:px-8'>
          <div className='w-full'>
            <Routes>
              <Route path='/' element={<Home />} />
              <Route path='/about' element={<About />} />
              <Route path='/allVideos' element={<AllVideos />} />
              <Route path='/admin' element={
                    <ProtectedRoute>
                        <AdminDashboard />
                    </ProtectedRoute>
                  }
              />
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
