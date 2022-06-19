pragma solidity ^0.4.24;

contract Greeter {
  address creator;
  string greeting;

  event BlockNumber(string str);
  constructor(string _greeting) public {
      creator = msg.sender;
      greeting = _greeting;
  }

  function greet() constant public returns (string) {return greeting; }

  function getBlockNumber() constant public returns (uint) {

     return block.number; 
      
  }

  function setGreeting(string _newgreeting) public { 
      greeting = _newgreeting;
      emit BlockNumber(_newgreeting);
   }
}