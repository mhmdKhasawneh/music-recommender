import './App.css'
import { BrowserRouter, Routes, Route } from "react-router-dom";
import Login from './pages/Login';
import Signup from './pages/Signup';
import Main from './pages/Main';

function App() {
  return (
    <div>
     <BrowserRouter>
         <Routes>
             <Route path='/' element={<Main/>}/>
             <Route path='/login' element={<Login/>}/>
             <Route path='/signup' element={<Signup/>}/>
         </Routes>
     </BrowserRouter>
    </div>
  )
}

export default App
