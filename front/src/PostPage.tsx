import { Post } from "./Post";
import "./PostPage.css";
import { Message } from "./Messsage";
import { Link, useParams } from "react-router-dom";
import { useEffect, useState } from "react";
import axios, { AxiosError } from "axios";
import { baseUrl } from "./const";

const client = axios.create({ withCredentials: true, timeout: 1000 });

const isAlllowShow: (p: Post) => boolean = (p: Post) => {
  if (p.limit_count && p.limit_count >= p.count) {
    return false
  }
  if(p.limit_date){
    const t = new Date();
    t.setTime(Date.parse(p.limit_date))
    if(t < new Date()){
      return false
    }
  }
  if(!p.require_message){
    return false
  }
  return true
}

const PostPage: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const parsed_id = parseInt(id);
  const initialPost: Post = {
    id: 0,
    user: "",
    email: "",
    title: "",
    kind: "",
    count: 0,
  };

  const initialPostMessage: Message & { view: string } = {
    post_id: parsed_id,
    user: "",
    email: "",
    message: "",
    view: "0",
  };

  const [post, setPost] = useState(initialPost);
  const [messages, setMessages] = useState([] as  Message[]);
  const [post_message, setPostMessage] = useState(initialPostMessage);
  const [url, seturl] = useState("");

  useEffect(() => {
    console.log("run useEffect");
    const p = async () => {
      console.log("start async get");
      const resp = await client
        .get(`${baseUrl}/api/post/${parsed_id}`)
        .catch((e: AxiosError<Response>) => {
          console.log(e);
          console.log(e.response);
        });
      console.log("end async get");
      console.log(resp);
      if (resp) {
        setPost(resp.data as Post);
        (resp.data as Post).url
          ? seturl((resp.data as Post).url!)
          : console.log("require message");
      }
    };
    p();
    const m = async () => {
      console.log("start async get");
      const resp = await client
        .get(`${baseUrl}/api/post/${parsed_id}/message`)
        .catch((e: AxiosError<Response>) => {
          console.log(e);
          console.log(e.response);
        });
      console.log("end async get");
      console.log(resp);
      if (resp) {
          setMessages(resp.data as Message[] || [] as Message[])
      }
    };
    m();
    console.log("finish use effect");
  }, [parsed_id]);

  const submit_message: (
    e: React.MouseEvent<HTMLInputElement>
  ) => Promise<void> = async (e) => {
    e.preventDefault();
    const url = `${baseUrl}/api/message`;
    const message_resp = await client
      .post(url, post_message, {
        headers: { "content-type": "application/json" },
      })
      .catch((e: AxiosError<Response>) => {
        console.log(e);
        console.log(e.response);
        alert("送信に失敗しました。\nしばらく時間をおいてお試しください。");
      });
    // catch 内でfunctionから抜けることができないのでここでrejectされているか確認する
    if (message_resp === undefined) {
      return;
    }
    console.log(message_resp);
    if (post_message.view === "1") {
      const post_resp = await client
        .get(`${baseUrl}/api/post/${id}`)
        .catch((e: AxiosError<Response>) => {
          console.log(e);
          console.log(e.response);
          alert("URLの取得に\nしばらく時間をおいてお試しください。");
        });
      // catch 内でfunctionから抜けることができないのでここでrejectされているか確認する
      if (post_resp === undefined) {
        return;
      }
      seturl(post_resp.data.url);
    }

    setMessages([
      ...messages,
      { ...post_message, inserted_at: new Date().toISOString() },
    ]);
  };
  return (
    <div>
      <h1>{post.title}</h1>
      <Link to="/">[戻る]</Link>
      <div className="post">
        <span className="user">{post.user}</span> [
        <a href={`mailto:${post.email}`}>E-mail</a>] {post.uploaded_at!}
        <br />
        {post.kind}
        {url !== "" ? `\n\n${url}` : undefined}
      </div>
      {messages.map((message) => (
        <>
          <hr />
          <div className="message" key={message.inserted_at}>
            <span className="user">{message.user}</span>[
            <a href={`mailto:${message.email}`}>E-mail</a>]{" "}
            {message.inserted_at}
            <br />
            {message.message}
          </div>
        </>
      ))}
      <hr />
      <form>
        追加発言
        <ul>
          <li>
            お名前
            <br />
            <input
              type="text"
              name="user"
              onChange={(e) =>
                setPostMessage({
                  ...post_message,
                  user: e.target.value,
                })
              }
            ></input>
          </li>
          <li>
            メールアドレス
            <br />
            <input
              type="text"
              name="email"
              onChange={(e) =>
                setPostMessage({
                  ...post_message,
                  email: e.target.value,
                })
              }
            ></input>
          </li>
          {
            isAlllowShow(post) ? (
          <li>
            URL表示
            <br />
            <input
              type="checkbox"
              name="view"
              value="1"
              onChange={(e) =>
                e.target.checked
                  ? setPostMessage({
                      ...post_message,
                      view: e.target.value,
                    })
                  : setPostMessage({ ...post_message, view: "" })
              }
            ></input>
            ここをチェックすることで提供品URLが表示されます。
          </li>
            ) : ""
          }
          <li>
            発言(エラー時再送信禁止)
            <br />
            <textarea
              name="message"
              rows={5}
              className="message-area"
              onChange={(e) =>
                setPostMessage({
                  ...post_message,
                  message: e.target.value,
                })
              }
            ></textarea>
          </li>
          <input type="submit" onClick={submit_message} value="送信"></input>
          <input type="reset" value="取消"></input>
        </ul>
      </form>
    </div>
  );
};

export default PostPage;
