import "./Menu.css";
import React from "react";
import { Link } from "react-router-dom";

const Menu: React.FC<{}> = () => {
  return (
    <nav>
      <ul>
        <li>back</li>
        <li>vip</li>
        <li>
          <Link to="/upload">upload</Link>
        </li>
        <li>log</li>
        <li>admin</li>
      </ul>
    </nav>
  );
};
export default Menu;
