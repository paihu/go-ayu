import axios, { AxiosError } from "axios";
import { useState } from "react";
import { Link, useHistory} from "react-router-dom";
import { baseUrl, Title } from "./const";
import { Post } from "./Post";
import {toast} from 'react-toastify';
const client = axios.create({ withCredentials: true, timeout: 1000 });

const initialPost: Post = {
  id: 0,
  user: "",
  email: "",
  title: "",
  kind: "",
  count: 0,
};

const Upload: React.FC<{}> = () => {
  const history = useHistory();
  const [post, cahngePost] = useState(initialPost);
  const submit_post: (e: React.MouseEvent<HTMLInputElement>) => Promise<void> =
    async (e) => {
      e.preventDefault();
      const url = `${baseUrl}/api/post`;
      const message_resp = await client
        .post(url, post, {
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
      toast("upload finish");
      //alert("upload finish")
      history.push(`/post/${message_resp.data.id}`);

    };
  return (
    <>
      <h1>{Title} 提供品UP</h1>
      <Link to="/">戻る</Link>
      <form>
        <hr />
        <ul>
          <li>
            お名前
            <br />
            <input
              type="text"
              onChange={(e) => cahngePost({ ...post, user: e.target.value })}
            ></input>
          </li>
          <li>
            メールアドレス
            <br />
            <input
              type="text"
              onChange={(e) => cahngePost({ ...post, email: e.target.value })}
            ></input>
          </li>
          <li>
            コメント
            <br />
            <textarea
              className="message-area"
              onChange={(e) => cahngePost({ ...post, comment: e.target.value })}
            ></textarea>
          </li>
          <li>
            削除用パスワード
            <br />
            <input
              type="password"
              onChange={(e) =>
                cahngePost({ ...post, delete_password: e.target.value })
              }
            ></input>
          </li>
        </ul>
        <hr />
        <ul>
          <li>
            提供品名
            <br />
            <input
              type="text"
              onChange={(e) => cahngePost({ ...post, title: e.target.value })}
            ></input>
          </li>
          <li>
            {" "}
            提供品種別
            <br />
            <input
              type="text"
              onChange={(e) => cahngePost({ ...post, kind: e.target.value })}
            ></input>
          </li>
          <li>
            提供品URL
            <br />
            <ul>
              <li>URLのみご記入ください</li>
            </ul>
            <textarea
              className="message-area"
              onChange={(e) => cahngePost({ ...post, url: e.target.value })}
            ></textarea>
          </li>
        </ul>
        <hr />
        <ul>
          <li>
            強制書き込み要求
            <br />
            <input
              type="checkbox"
              onChange={(e) =>
                e.target.checked
                  ? cahngePost({ ...post, require_message: true })
                  : cahngePost({ ...post, require_message: false })
              }
            ></input>
            全員に書き込みをさせる
          </li>
          <li>
            ダウンロード数制限
            <br />
            <input
              type="text"
              onChange={(e) =>
                cahngePost({ ...post, limit_count: parseInt(e.target.value) })
              }
            ></input>
            (制限なしの場合は「0」)
          </li>
        </ul>
        <hr />
        <ul>
          <input type="submit" value="送信" onClick={submit_post}></input>
          <input type="reset" value="取消"></input>
        </ul>
      </form>
    </>
  );
};

export default Upload;
