import './App.css';
import React, {useState,useEffect} from 'react';
import { createTheme, ThemeProvider} from '@mui/material/styles';
import * as axios from 'axios';
import { RoomView } from './RoomView';
import {Paper, CssBaseline,ListItem} from  "@mui/material";


function App() {

  const theme = createTheme({
    shape:{
      borderRadius: 0
    },
    palette:{
      mode:"dark"
    }
  });

  const [msg,setMsg] = useState([]);
  const [text,setText] = useState("");
  const [counter,setCounter] = useState(0);



  useEffect(()=>getData(),[counter])

  async function getData(){
    let response = await axios.get("http://127.0.0.1:8081/",{headers:{'Accept': 'application/json'}});
    console.log(response)
    if (response.status === 200){
      let data = response.data
      if (data != null){
       setMsg(prevMsg => [...prevMsg,...data])
      }
    }
    const interval = await new Promise(r=> setTimeout(r,10000));
    setCounter(prev=>prev+1) 
    return ()=>clearInterval(interval)
    
  };

  function handleText({target}){
    setText(target.value)
  };
  
  async function handleSubmit(index){
    let body = ""
    setMsg(prev => [...prev, ("client: "+text)])
    switch (index){
      case 0://send message
        body = "/msg " + text;      break;
      case 1://send join room
        body = "/join " + text;     break;
      case 2://send change nickname
        body = "/nick " + text;     break;
      case 3://create room
        body = "/create " + text;   break;
      case 4:
        body = "/rooms";            break;
      default:
        break;
    }
    console.log(body)
    let response = await axios({
      method: 'post',
      url: 'http://127.0.0.1:8081/',
      data: body
    });
    if (response.status === 200) {
      setMsg(prevMsg => [...prevMsg])
      setText('')
    }

  }

  const listItems = msg.map((value,index) =>{
    if (value.startsWith("client: ")){
      let temp = value.slice(7)
        return <ListItem key = {index} sx={{display:"inline-flex",justifyContent:"flex-end"}}>{temp+ " <"}</ListItem>
    }
    return <ListItem key = {index}>{value}</ListItem>
  });



  return (
    <div className="App">
      <ThemeProvider theme={theme}  >
        <CssBaseline/>
        <Paper sx={{display:"flex", backgroundColor:"red", height:"100vh",flexDirection:"coloumn"}}>
          <RoomView text = {text} handleSubmit={handleSubmit} handleText={handleText} listItems={listItems}  />
        </Paper>
      </ThemeProvider>
    </div>
  );
}



export default App;
