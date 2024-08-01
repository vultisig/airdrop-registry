import Image from "next/image";

export default function Sections() {
    return (
        <div className="min-h-screen flex flex-col items-center p-4">
        {/* <!-- Header --> */}
        <div className="w-full flex items-center justify-between p-4 bg-[#061B3A] rounded-lg">
            <div className="text-xl font-semibold flex items-center">
                <Image src="/img/logo.png" width={48} height={48} alt="logo" />
                <h1 className="text-[40px] font-semibold">Vultisig</h1>
            </div>
            <div className="flex items-center space-x-4">
                <button className="rounded-full ">
                    <span className="sr-only">Profile</span>
                    {/* <!-- Icon for profile --> */}
                    <Image src="/img/UserCircle.png" width={35} height={35} alt="user" />
                </button>
            </div>
        </div>
        {/* <!-- Main Content --> */}
        <div className="w-full max-w-4xl mt-8">
            <div className="text-start mb-8">
                {/* <select className="p-2 bg-[#061B3A] rounded-lg px-10 py-2">
                    <option>Main Vault</option>
                </select> */}

                <button id="dropdownDefaultButton" data-dropdown-toggle="dropdown" className="text-white bg-[#061B3A]  focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded-lg text-sm px-20 py-2.5 text-left inline-flex items-center  dark:focus:ring-blue-800" type="button">
                    Main Vault
                    <svg className="w-2.5 h-2.5 ms-3" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 10 6">
                        <path stroke="currentColor" strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="m1 1 4 4 4-4" />
                    </svg>
                </button>

                {/* <!-- Dropdown menu --> */}
                <div id="dropdown" className="z-10 hidden bg-white divide-y divide-gray-100 rounded-lg shadow w-44 dark:bg-gray-700">
                    <ul className="py-2 text-sm text-gray-700 dark:text-gray-200" aria-labelledby="dropdownDefaultButton">
                        <li>
                            <a href="#" className="block px-4 py-2 hover:bg-gray-100 dark:hover:bg-gray-600 dark:hover:text-white">Dashboard</a>
                        </li>
                        <li>
                            <a href="#" className="block px-4 py-2 hover:bg-gray-100 dark:hover:bg-gray-600 dark:hover:text-white">Settings</a>
                        </li>
                        <li>
                            <a href="#" className="block px-4 py-2 hover:bg-gray-100 dark:hover:bg-gray-600 dark:hover:text-white">Earnings</a>
                        </li>
                        <li>
                            <a href="#" className="block px-4 py-2 hover:bg-gray-100 dark:hover:bg-gray-600 dark:hover:text-white">Sign out</a>
                        </li>
                    </ul>
                </div>

            </div>
            <div>
                <p>
                    Total Balance
                </p>
            </div>
            <div className="text-3xl font-bold text-start mb-8">$365,899.00</div>
            {/* <!-- Asset List --> */}
            <div className="space-y-4">
                {/* <!-- Bitcoin --> */}
                <div className="bg-[#061B3A] p-4 rounded-lg flex items-center justify-between flex-wrap">
                    <div className="flex items-center space-x-4">
                        <img src="/img/Bitcoin.png" alt="Bitcoin" className="h-8 w-8" />
                        <div>
                            <div className="font-semibold">Bitcoin</div>
                            <div>BTC</div>

                        </div>
                    </div>
                    <div className="text-sm text-[#81F8DE]">bc1prsjtwm7682v6nhx2uwfgcferennd7pcvcqv7v6w</div>
                    <div className="text-right flex items-center">
                        <div className="font-semibold text-xs text-[#F3F4F5] mx-3">1.1 BTC</div>
                        <div className="text-xl">$65,899</div>
                    </div>
                    <div className="flex items-center space-x-4">
                        <Image src="/img/Cube.png" width={22} height={22} alt="cube" />
                        <Image src="/img/QR.png" width={22} height={22} alt="cube" />
                        <Image src="/img/COPY.png" width={22} height={22} alt="cube" />
                        <Image src="/img/CaretLeft.png" width={22} height={22} alt="cube" />

                    </div>
                </div>
                {/* <!-- Ethereum --> */}
                <div className="bg-[#061B3A] p-4 rounded-lg flex items-center justify-between flex-wrap">
                    <div className="flex items-center space-x-4">
                        <img src="/img/et.png" alt="Ethereum" className="h-8 w-8 rounded-full" />
                        <div>
                            <div className="font-semibold">Ethereum</div>
                            <div>BTC</div>
                        </div>
                    </div>
                    <div className="text-sm text-[#81F8DE]">0xcB1dA24292b889862f599AC5B10F42b06DE07e4</div>

                    <div className="text-right flex items-center">
                        <div className="font-semibold mx-2 bg-[#11284A] p-1 text-xs text-center rounded-full">3 assets</div>
                        <div className="text-xl">$65,899</div>
                    </div>
                    <div className="flex items-center space-x-4">
                        <Image src="/img/Cube.png" width={22} height={22} alt="cube" />
                        <Image src="/img/QR.png" width={22} height={22} alt="cube" />
                        <Image src="/img/COPY.png" width={22} height={22} alt="cube" />
                        <Image src="/img/CaretLeft.png" width={22} height={22} alt="cube" />

                    </div>
                </div>
                {/* <!-- Solana --> */}
                <div className="bg-[#061B3A] p-4 rounded-lg flex items-center justify-between flex-wrap">
                    <div className="flex items-center space-x-4">
                        <img src="/img/Solana.png" alt="Solana" className="h-8 w-8 rounded-full" />
                        <div>
                            <div className="font-semibold">Solana</div>
                            <div>BTC</div>
                        </div>
                    </div>
                    <div className="text-sm text-[#81F8DE]">ELPecyZbKieS2NuNAGPzna6q7r8DYa7aFaptDo7X1GMJ</div>

                    <div className="text-right flex items-center">
                        <div className="font-semibold mx-2 bg-[#11284A] p-1 text-xs text-center rounded-full">3 assets</div>
                        <div className="text-xl">$65,899</div>
                    </div>
                    <div className="flex items-center space-x-4">
                        <Image src="/img/Cube.png" width={22} height={22} alt="cube" />
                        <Image src="/img/QR.png" width={22} height={22} alt="cube" />
                        <Image src="/img/COPY.png" width={22} height={22} alt="cube" />
                        <Image src="/img/CaretLeft.png" width={22} height={22} alt="cube" />

                    </div>
                </div>
                {/* <!-- THORChain --> */}
                <div className="bg-[#061B3A] p-4 rounded-lg flex items-center justify-between flex-wrap">
                    <div className="flex items-center space-x-4">
                        <img src="/img/Thorchain.png" alt="THORChain" className="h-8 w-8 rounded-full" />
                        <div>
                            <div className="font-semibold">THORChain</div>
                            <div>BTC</div>
                        </div>
                    </div>
                    <div className="text-sm text-[#81F8DE]">thor1cfelrennd7pcvcqv7v6w7682v6nhx2uwfg</div>

                    <div className="text-right flex items-center">
                        <div className="font-semibold text-xs text-[#F3F4F5] mx-3">1.1 BTC</div>
                        <div className="text-xl">$65,899</div>
                    </div>
                    <div className="flex items-center space-x-4">
                        <Image src="/img/Cube.png" width={22} height={22} alt="cube" />
                        <Image src="/img/QR.png" width={22} height={22} alt="cube" />
                        <Image src="/img/COPY.png" width={22} height={22} alt="cube" />
                        <Image src="/img/CaretLeft.png" width={22} height={22} alt="cube" />
                    </div>
                </div>
                {/* <!-- Add New Chain --> */}
                <div className="text-left">

                    <button className="rounded-full flex items-center">

                        {/* <!-- Icon for add new chain --> */}
                        <Image src="/img/PlusCircle.png" width={40} height={40} alt="PlusCircle" />

                        <span className="mx-2 text-base">Choose Chains</span>
                    </button>
                </div>
            </div>
        </div>
    </div>
    );
}