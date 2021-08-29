import axios from "axios";
import { useEffect, useState } from "react";
import { Link, useLocation } from "react-router-dom";
import { baseUrl } from "./const";
import "./Content.css";
import { Post } from "./Post";

const client = axios.create({ withCredentials: true, timeout: 1000 });

interface ContentProps {
  page: number;
}

const delete_post: (p: Post) => () => void = (post: Post) => {
  return () => {
    console.log(post.id);
  };
};
const parseDate: (s: string) => string = (s: string) => {
  const t = new Date();
  t.setTime(Date.parse(s));
  return `${t.getFullYear()}/${(t.getMonth() + 1)
    .toString()
    .padStart(2, "0")}/${t.getDate().toString().padStart(2, "0")}`;
};
const Content: React.FC<ContentProps> = ({ page }) => {
  const location = useLocation();
  const [posts, setPosts] = useState([] as Post[]);
  useEffect(() => {
    const f = async () => {
      const resp = await client
        .get(`${baseUrl}/api/post/page/${page}`)
        .catch((e) => {
          console.log(e);
        });
      resp ? setPosts(resp.data as Post[]) : console.log(resp);
    };
    f();
  }, [page, location]);
  return (
    <form>
      <table>
        <thead>
          <tr>
            <th>
              削<br />除
            </th>
            <th>提供日</th>
            <th>UP職人</th>
            <th>提供品</th>
            <th>種別</th>
            <th>DL数</th>
            <th>トピック状況</th>
          </tr>
        </thead>
        <tbody>
          {posts.map((post) => (
            <tr>
              <td onClick={delete_post(post)}>
                <input type="checkbox" name={`${post.id}`} value="1" />
              </td>
              <td>{parseDate(post.uploaded_at!)}</td>
              <td>{post.user}</td>
              <td>
                <Link to={`/post/${post.id}`}>{post.title}</Link>
              </td>
              <td>{post.kind}</td>
              <td>{post.count ? post.count : 0}</td>
              <td>
                {post.limit_count
                  ? `あと${post.limit_count - (post.count ? post.count : 0)}人 `
                  : ""}
                {post.require_message ? "要書込 " : ""}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
      <div>
        <input type="password"></input>
        <input type="submit" value="削除" />
      </div>
    </form>
  );
};

export default Content;
