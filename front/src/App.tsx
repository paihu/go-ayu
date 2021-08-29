import "./App.css";
import Content from "./Content";
import Menu from "./Menu";
import Notice from "./Notice";
import PostPage from "./PostPage";
import { BrowserRouter as Router, Switch, Route } from "react-router-dom";
import { Title } from "./const";
import Upload from "./Upload";
import {ToastContainer} from "react-toastify";
import 'react-toastify/dist/ReactToastify.min.css';

const App: React.FC<{}> = () => {
  return (
    <Router>
      <ToastContainer/>
      <Switch>
        <Route exact path="/">
          <div className="App">
            <h1>{Title}</h1>
            <Menu></Menu>
            <hr />
            <Notice></Notice>
            <hr />
            <Content page={0}></Content>
          </div>
        </Route>
        <Route exact path="/post/:id">
          <PostPage></PostPage>
        </Route>
        <Route exact path="/upload">
          <Upload></Upload>
        </Route>
      </Switch>
    </Router>
  );
};

export default App;
