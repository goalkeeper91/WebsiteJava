import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import Header from './components/Header';
import Footer from './components/Footer';
import Home from './pages/Home';
import About from './pages/About';
import AllVideos from './pages/AllVideos';

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
            </Routes>
          </div>
        </main>

        <Footer />
      </div>
    </Router>
  );
};

export default App;
