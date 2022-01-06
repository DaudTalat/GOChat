import React from 'react';
import './App.css';
import {TextField, Paper, List} from  "@mui/material";
import SplitButton from './SplitButton';


export const RoomView = (props)=>{
    return (
          <Paper elevation={1} sx = {{height:"100vh",float:"left", width:"100vw"}}>
            <Paper> </Paper>
            <Paper elevation={0} sx={{height:"94vh", padding:"2vh 0 2vh", overflow:"auto",display:"flex",flexDirection:"column-reverse"}}>
              <List sx = {{}}>
                {props.listItems}
              </List>
            </Paper>
            <Paper elevation={1} sx ={{display:"flex",position:"absolute",bottom:"0", width:"100vw"}}>
              <TextField 
                multiline 
                sx = {{flex:"9", margin:"0"}}
                value = {props.text}
                onChange = {props.handleText}
              ></TextField>
              <SplitButton handleSubmit={props.handleSubmit} ></SplitButton>
            </Paper>
          </Paper>
    );
  }