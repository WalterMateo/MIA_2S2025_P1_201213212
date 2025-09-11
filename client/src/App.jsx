import { useState } from 'react'
import reactLogo from './assets/react.svg'
import viteLogo from '/vite.svg'
import './App.css'
import { BrowserRouter, Route, Routes } from 'react-router-dom'
import { Contenedores } from './components/Contenedores.jsx'
import ReactDOM from 'react-dom';
import React from 'react';

function App() {
  const [count, setCount] = useState(0)

  return (
     <BrowserRouter>
          <Routes>
            <Route path="" element={<Contenedores />} />
            
          </Routes>
        </BrowserRouter>  
  )
}

export default App
