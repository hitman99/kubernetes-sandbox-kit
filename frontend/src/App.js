import logo from './logo.svg';
import './App.css';
import Registration from './registration';
import React, {useState} from "react";
import {loadState} from "./utils/storage";
import Readme from './readme';

function App() {
  const [regData, setRegData] = useState(loadState());
  console.log(regData);
  let view;
  if( !regData.user.id ) {
    view = <Registration regData={regData} setRegData={setRegData}/>
  } else {
    view = <Readme regData={regData} setRegData={setRegData}/>
  }
  return (
    <div className="App">
      { view }
    </div>
  );
}

export default App;
