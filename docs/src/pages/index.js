import React, { useEffect } from 'react';

export default function Home() {
  useEffect(function onFirstMount() {
    window.location.href='https://dataconnector.app';
  }, []); 
  
  return null;
};
