import "./Notice.css";
import React from "react";

const Notice: React.FC<{}> = () => {
  return (
    <div className="Notice">
      <ul>
        <li className="warn">
          要書込みのものは追加書き込みをすることでアドレスが現れます。
        </li>
        <li>
          UP職人はDLできる人数を設定できます。設定人数を超えると提供品のアドレスは表示されません。
        </li>
        <li>優秀UP職人として活動していただいた方にはVIP特典が贈られます。</li>
      </ul>
    </div>
  );
};

export default Notice;
